from dslabmp import Context, Message, Process

TIMER_DELAY = 100

# AT MOST ONCE ---------------------------------------------------------------------------------------------------------

class AtMostOnceSender(Process):
    def __init__(self, proc_id: str, receiver_id: str):
        self._id = proc_id
        self._receiver = receiver_id
        self._index = 0

    def on_local_message(self, msg: Message, ctx: Context):
        # receive message for delivery from local user
        self._index += 1
        msg._data['id'] = self._index
        ctx.send(msg, self._receiver)

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver here
        pass

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        pass


class AtMostOnceReceiver(Process):
    def __init__(self, proc_id: str):
        self._id = proc_id
        self._buffer = 0

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver
        # deliver message to local user with ctx.send_local()
        message_id = msg._data['id']
        if not self._buffer & (1 << message_id):
            self._buffer |= (1 << message_id)
            del msg._data['id']
            ctx.send_local(msg)

    def on_local_message(self, msg: Message, ctx: Context):
        # not used in this task
        pass

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        pass


# AT LEAST ONCE --------------------------------------------------------------------------------------------------------

class AtLeastOnceSender(Process):
    def __init__(self, proc_id: str, receiver_id: str):
        self._receiver = receiver_id
        self._id = proc_id
        self._index = 0
        self._sent = {}

    def on_local_message(self, msg: Message, ctx: Context):
        # receive message for delivery from local user
        self._sent[self._index] = msg._data['text']
        msg._data['from'] = self._id
        msg._data['id'] = self._index
        ctx.send(msg, self._receiver)
        ctx.set_timer(str(self._index), TIMER_DELAY)
        self._index += 1

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver here
        message_id = msg._data['id']
        if message_id in self._sent:
            del self._sent[message_id]

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        message_index = int(timer_name)
        if message_index in self._sent:
            ctx.send(Message('MESSAGE', {
                'text': self._sent[message_index],
                'id': message_index,
                'from': self._id
            }), self._receiver)
            ctx.set_timer(timer_name, TIMER_DELAY)


class AtLeastOnceReceiver(Process):
    def __init__(self, proc_id: str):
        self._id = proc_id

    def on_local_message(self, msg: Message, ctx: Context):
        # not used in this task
        pass

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver
        # deliver message to local user with ctx.send_local()
        response_msg = Message('MESSAGE', {'id': msg._data['id']})
        ctx.send(response_msg, msg._data['from'])

        # Cleaning up unnecessary data before sending locally
        del msg._data['id']
        del msg._data['from']

        ctx.send_local(msg)

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        pass


# EXACTLY ONCE ---------------------------------------------------------------------------------------------------------

class ExactlyOnceSender(Process):
    def __init__(self, proc_id: str, receiver_id: str):
        self._id = proc_id
        self._receiver = receiver_id
        self._sent = {}
        self._index = 0

    def on_local_message(self, msg: Message, ctx: Context):
        # receive message for delivery from local user
        msg._data['from'] = self._id
        msg._data['id'] = self._index
        self._sent[self._index] = msg._data['text']
        ctx.send(msg, self._receiver)
        ctx.set_timer(str(self._index), TIMER_DELAY)
        self._index += 1

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver here
        message_id = msg._data['id']
        if message_id in self._sent:
            del self._sent[message_id]

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        message_index = int(timer_name)
        if message_index in self._sent:
            message_content = {
                'text': self._sent[message_index],
                'id': message_index,
                'from': self._id
            }
            ctx.send(Message('MESSAGE', message_content), self._receiver)
            ctx.set_timer(timer_name, TIMER_DELAY)


class ExactlyOnceReceiver(Process):
    def __init__(self, proc_id: str):
        self._id = proc_id
        self._buffer = 0

    def on_local_message(self, msg: Message, ctx: Context):
        # not used in this task
        pass

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver
        # deliver message to local user with ctx.send_local()
        message_id = msg._data['id']
        ctx.send(Message('MESSAGE', {'id': message_id}), msg._data['from'])

        if not self._buffer & (1 << message_id):
            self._buffer |= (1 << message_id)
            del msg._data['id']
            del msg._data['from']
            ctx.send_local(msg)

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        pass


# EXACTLY ONCE + ORDERED -----------------------------------------------------------------------------------------------

class ExactlyOnceOrderedSender(Process):
    def __init__(self, proc_id: str, receiver_id: str):
        self._id = proc_id
        self._receiver = receiver_id
        self._queue = []
        self._sent = {}
        self._index = 0

    def on_local_message(self, msg: Message, ctx: Context):
        # receive message for delivery from local user
        if not self._sent:
            message_content = {
                'text': msg._data['text'],
                'id': self._index,
                'from': self._id
            }
            ctx.send(Message('MESSAGE', message_content), self._receiver)
            ctx.set_timer(str(self._index), TIMER_DELAY)
            self._sent[self._index] = msg._data['text']
            self._index += 1
        else:
            self._queue.append(msg._data['text'])

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver here
        message_id = msg._data['id']
        if message_id in self._sent:
            del self._sent[message_id]

        if not self._sent and self._queue:
            message_content = {
                'text': self._queue.pop(0),
                'id': self._index,
                'from': self._id
            }
            ctx.send(Message('MESSAGE', message_content), self._receiver)
            ctx.set_timer(str(self._index), TIMER_DELAY)
            self._sent[self._index] = message_content['text']
            self._index += 1

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        message_index = int(timer_name)
        if message_index in self._sent:
            message_content = {
                'text': self._sent[message_index],
                'id': message_index,
                'from': self._id
            }
            ctx.send(Message('MESSAGE', message_content), self._receiver)
            ctx.set_timer(timer_name, TIMER_DELAY)


class ExactlyOnceOrderedReceiver(Process):
    def __init__(self, proc_id: str):
        self._id = proc_id
        self._buffer = 0

    def on_local_message(self, msg: Message, ctx: Context):
        # not used in this task
        pass

    def on_message(self, msg: Message, sender: str, ctx: Context):
        # process messages from receiver
        # deliver message to local user with ctx.send_local()
        message_id = msg._data['id']
        ctx.send(Message('MESSAGE', {'id': message_id}), msg._data['from'])

        if not self._buffer & (1 << message_id):
            self._buffer |= (1 << message_id)
            del msg._data['id']
            del msg._data['from']
            ctx.send_local(msg)

    def on_timer(self, timer_name: str, ctx: Context):
        # process fired timers here
        pass
