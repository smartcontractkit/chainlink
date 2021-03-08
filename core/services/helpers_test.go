package services

func (ht *HeadTracker) ExportedDone() chan struct{} {
	return ht.done
}
