class WorkerProfile:
    __slots__ = ['worker_id', 'wage', 'hours', 'pto_used']
    def __init__(self, worker_id, wage, hours=40, pto_used=2):
        self.worker_id = worker_id
        self.wage = wage
        self.hours = hours
        self.pto_used = pto_used
    def net_amount(self):
        return (self.hours + self.pto_used) * self.wage
