from dslabmp import Context, Message, Process
from typing import List, Tuple, Set

class BroadcastProcess(Process):
    def __init__(self, proc_id: str, processes: List[str]):
        self._id = proc_id
        self._processes = processes
        self._delivered = set()
        self._acks = []

    def on_local_message(self, msg: Message, ctx: Context):
        if msg.type == 'SEND':
            for proc in self._processes:
                ctx.send(Message('PREPARE', {'text': msg['text']}), proc)
            self._acks.append([msg['text'], [self._id], 1])
            ctx.set_timer('RETRY_{}'.format(msg['text']), 10.0)

    def on_message(self, msg: Message, sender: str, ctx: Context):
        if msg.type == 'PREPARE':
            idx, count = self._update_acks(msg['text'], sender)
            ctx.send(Message('ACK', {'text': msg['text']}), sender)
            if count > len(self._processes) // 2 and msg['text'] not in self._delivered:
                for proc in self._processes:
                    ctx.send(Message('COMMIT', {'text': msg['text']}), proc)

        elif msg.type == 'ACK':
            idx, count = self._update_acks(msg['text'], sender)
            if count > len(self._processes) // 2 and msg['text'] not in self._delivered:
                for proc in self._processes:
                    ctx.send(Message('COMMIT', {'text': msg['text']}), proc)

        elif msg.type == 'COMMIT':
            if msg['text'] not in self._delivered:
                self._delivered.add(msg['text'])
                ctx.send_local(Message('DELIVER', {'text': msg['text']}))
                ctx.cancel_timer('RETRY_{}'.format(msg['text']))

    def on_timer(self, timer_name: str, ctx: Context):
        if timer_name.startswith('RETRY_'):
            text = timer_name.split('_')[1]
            for proc in self._processes:
                ctx.send(Message('PREPARE', {'text': text}), proc)

    def _update_acks(self, text: str, sender: str) -> Tuple[int, int]:
        for idx, (msg_text, senders, count) in enumerate(self._acks):
            if msg_text == text:
                if sender not in senders:
                    senders.append(sender)
                    count += 1
                self._acks[idx] = (msg_text, senders, count)
                return idx, count

        self._acks.append([text, [self._id, sender], 2])
        return len(self._acks) - 1, 2
