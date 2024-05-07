package workflow

import (
    "errors"
    "fmt"

    "github.com/smartcontractkit/chainlink-common/pkg/values"

    "github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)



type mergeRunner2[ I,II, O any] struct {
	mergeRunnerBase
	fn func(I, II, ) (O, error)
}

func (m mergeRunner2[ I,II, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner2[any,any, any]{}

func Merge2[ I,II, O any](ref string, wb1 *Builder[I], wb2 *Builder[II],  merge func(I,II,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false

    mr := &mergeRunner2[ I,II, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner3[ I,II,III, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, ) (O, error)
}

func (m mergeRunner3[ I,II,III, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner3[any,any,any, any]{}

func Merge3[ I,II,III, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III],  merge func(I,II,III,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false

    mr := &mergeRunner3[ I,II,III, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner4[ I,II,III,IIII, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, IIII, ) (O, error)
}

func (m mergeRunner4[ I,II,III,IIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
    v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,v4,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner4[any,any,any,any, any]{}

func Merge4[ I,II,III,IIII, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III], wb4 *Builder[IIII],  merge func(I,II,III,IIII,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb4.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false
	wb1.root.open[wb4.current.Ref()] = false

    mr := &mergeRunner4[ I,II,III,IIII, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,wb4.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner5[ I,II,III,IIII,IIIII, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, IIII, IIIII, ) (O, error)
}

func (m mergeRunner5[ I,II,III,IIII,IIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
    v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
    v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,v4,v5,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner5[any,any,any,any,any, any]{}

func Merge5[ I,II,III,IIII,IIIII, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III], wb4 *Builder[IIII], wb5 *Builder[IIIII],  merge func(I,II,III,IIII,IIIII,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb4.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb5.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false
	wb1.root.open[wb4.current.Ref()] = false
	wb1.root.open[wb5.current.Ref()] = false

    mr := &mergeRunner5[ I,II,III,IIII,IIIII, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,wb4.current,wb5.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner6[ I,II,III,IIII,IIIII,IIIIII, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, IIII, IIIII, IIIIII, ) (O, error)
}

func (m mergeRunner6[ I,II,III,IIII,IIIII,IIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
    v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
    v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
    v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,v4,v5,v6,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner6[any,any,any,any,any,any, any]{}

func Merge6[ I,II,III,IIII,IIIII,IIIIII, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III], wb4 *Builder[IIII], wb5 *Builder[IIIII], wb6 *Builder[IIIIII],  merge func(I,II,III,IIII,IIIII,IIIIII,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb4.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb5.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb6.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false
	wb1.root.open[wb4.current.Ref()] = false
	wb1.root.open[wb5.current.Ref()] = false
	wb1.root.open[wb6.current.Ref()] = false

    mr := &mergeRunner6[ I,II,III,IIII,IIIII,IIIIII, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,wb4.current,wb5.current,wb6.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner7[ I,II,III,IIII,IIIII,IIIIII,IIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, IIII, IIIII, IIIIII, IIIIIII, ) (O, error)
}

func (m mergeRunner7[ I,II,III,IIII,IIIII,IIIIII,IIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
    v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
    v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
    v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
    v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,v4,v5,v6,v7,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner7[any,any,any,any,any,any,any, any]{}

func Merge7[ I,II,III,IIII,IIIII,IIIIII,IIIIIII, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III], wb4 *Builder[IIII], wb5 *Builder[IIIII], wb6 *Builder[IIIIII], wb7 *Builder[IIIIIII],  merge func(I,II,III,IIII,IIIII,IIIIII,IIIIIII,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb4.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb5.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb6.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb7.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false
	wb1.root.open[wb4.current.Ref()] = false
	wb1.root.open[wb5.current.Ref()] = false
	wb1.root.open[wb6.current.Ref()] = false
	wb1.root.open[wb7.current.Ref()] = false

    mr := &mergeRunner7[ I,II,III,IIII,IIIII,IIIIII,IIIIIII, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,wb4.current,wb5.current,wb6.current,wb7.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner8[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, ) (O, error)
}

func (m mergeRunner8[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
    v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
    v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
    v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
    v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}
    v8, err := capabilities.UnwrapValue[IIIIIIII](ls[8-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,v4,v5,v6,v7,v8,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner8[any,any,any,any,any,any,any,any, any]{}

func Merge8[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III], wb4 *Builder[IIII], wb5 *Builder[IIIII], wb6 *Builder[IIIIII], wb7 *Builder[IIIIIII], wb8 *Builder[IIIIIIII],  merge func(I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb4.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb5.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb6.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb7.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb8.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false
	wb1.root.open[wb4.current.Ref()] = false
	wb1.root.open[wb5.current.Ref()] = false
	wb1.root.open[wb6.current.Ref()] = false
	wb1.root.open[wb7.current.Ref()] = false
	wb1.root.open[wb8.current.Ref()] = false

    mr := &mergeRunner8[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,wb4.current,wb5.current,wb6.current,wb7.current,wb8.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner9[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, IIIIIIIII, ) (O, error)
}

func (m mergeRunner9[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
    v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
    v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
    v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
    v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}
    v8, err := capabilities.UnwrapValue[IIIIIIII](ls[8-1])
	if err != nil {
		return nil, false, err
	}
    v9, err := capabilities.UnwrapValue[IIIIIIIII](ls[9-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,v4,v5,v6,v7,v8,v9,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner9[any,any,any,any,any,any,any,any,any, any]{}

func Merge9[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III], wb4 *Builder[IIII], wb5 *Builder[IIIII], wb6 *Builder[IIIIII], wb7 *Builder[IIIIIII], wb8 *Builder[IIIIIIII], wb9 *Builder[IIIIIIIII],  merge func(I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb4.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb5.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb6.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb7.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb8.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb9.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false
	wb1.root.open[wb4.current.Ref()] = false
	wb1.root.open[wb5.current.Ref()] = false
	wb1.root.open[wb6.current.Ref()] = false
	wb1.root.open[wb7.current.Ref()] = false
	wb1.root.open[wb8.current.Ref()] = false
	wb1.root.open[wb9.current.Ref()] = false

    mr := &mergeRunner9[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,wb4.current,wb5.current,wb6.current,wb7.current,wb8.current,wb9.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


type mergeRunner10[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII,IIIIIIIIII, O any] struct {
	mergeRunnerBase
	fn func(I, II, III, IIII, IIIII, IIIIII, IIIIIII, IIIIIIII, IIIIIIIII, IIIIIIIIII, ) (O, error)
}

func (m mergeRunner10[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII,IIIIIIIIII, O]) Run(value values.Value) (values.Value, bool, error) {
	ls := value.(*values.List).Underlying
    v1, err := capabilities.UnwrapValue[I](ls[1-1])
	if err != nil {
		return nil, false, err
	}
    v2, err := capabilities.UnwrapValue[II](ls[2-1])
	if err != nil {
		return nil, false, err
	}
    v3, err := capabilities.UnwrapValue[III](ls[3-1])
	if err != nil {
		return nil, false, err
	}
    v4, err := capabilities.UnwrapValue[IIII](ls[4-1])
	if err != nil {
		return nil, false, err
	}
    v5, err := capabilities.UnwrapValue[IIIII](ls[5-1])
	if err != nil {
		return nil, false, err
	}
    v6, err := capabilities.UnwrapValue[IIIIII](ls[6-1])
	if err != nil {
		return nil, false, err
	}
    v7, err := capabilities.UnwrapValue[IIIIIII](ls[7-1])
	if err != nil {
		return nil, false, err
	}
    v8, err := capabilities.UnwrapValue[IIIIIIII](ls[8-1])
	if err != nil {
		return nil, false, err
	}
    v9, err := capabilities.UnwrapValue[IIIIIIIII](ls[9-1])
	if err != nil {
		return nil, false, err
	}
    v10, err := capabilities.UnwrapValue[IIIIIIIIII](ls[10-1])
	if err != nil {
		return nil, false, err
	}

	merged, err :=  m.fn(v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(merged)
	return wrapped, true, err
}

var _ capability = &mergeRunner10[any,any,any,any,any,any,any,any,any,any, any]{}

func Merge10[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII,IIIIIIIIII, O any](ref string, wb1 *Builder[I], wb2 *Builder[II], wb3 *Builder[III], wb4 *Builder[IIII], wb5 *Builder[IIIII], wb6 *Builder[IIIIII], wb7 *Builder[IIIIIII], wb8 *Builder[IIIIIIII], wb9 *Builder[IIIIIIIII], wb10 *Builder[IIIIIIIIII],  merge func(I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII,IIIIIIIIII,) (O, error)) (*Builder[O], error) {
        if wb1.root != wb2.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb3.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb4.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb5.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb6.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb7.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb8.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb9.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }
        if wb1.root != wb10.root {
            return nil, errors.New("cannot merge builders from different workflows")
        }

    wb1.root.lock.Lock()
	defer wb1.root.lock.Unlock()

	if wb1.root.names[ref] {
		return nil, fmt.Errorf("name %s already exists as a step", ref)
	}
	wb1.root.names[ref] = true
	wb1.root.open[ref] = true
	wb1.root.open[wb1.current.Ref()] = false
	wb1.root.open[wb2.current.Ref()] = false
	wb1.root.open[wb3.current.Ref()] = false
	wb1.root.open[wb4.current.Ref()] = false
	wb1.root.open[wb5.current.Ref()] = false
	wb1.root.open[wb6.current.Ref()] = false
	wb1.root.open[wb7.current.Ref()] = false
	wb1.root.open[wb8.current.Ref()] = false
	wb1.root.open[wb9.current.Ref()] = false
	wb1.root.open[wb10.current.Ref()] = false

    mr := &mergeRunner10[ I,II,III,IIII,IIIII,IIIIII,IIIIIII,IIIIIIII,IIIIIIIII,IIIIIIIIII, O]{
        fn: merge,
        mergeRunnerBase: mergeRunnerBase{
        nonTriggerCapability{
            inputs: mergeOutputs(wb1.current,wb2.current,wb3.current,wb4.current,wb5.current,wb6.current,wb7.current,wb8.current,wb9.current,wb10.current,),
            ref:    ref,
        },
      },
    }
    wb1.root.spec.Actions = append(wb1.root.spec.Actions, capabilityToStepDef(mr))
    wb1.root.capabilities = append(wb1.root.capabilities, mr)
    wb1.root.spec.LocalExecutions[ref] = mr
    return &Builder[O]{
        root: wb1.root,
        current: mr,
    }, nil
}


