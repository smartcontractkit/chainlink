package edwards25519

// AllowVarTime sets a flag in this object which determines if a faster
// but variable time implementation can be used. Set this only on Points
// which represent public information. Using variable time algorithms to
// operate on private information can result in timing side-channels.
func (P *point) AllowVarTime(varTime bool) {
	P.varTime = varTime
}
