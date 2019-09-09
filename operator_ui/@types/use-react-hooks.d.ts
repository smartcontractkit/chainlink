declare module 'use-react-hooks' {
  import { FC } from 'react'

  export {
    useCallback,
    useContext,
    useDebugValue,
    useEffect,
    useImperativeHandle,
    useLayoutEffect,
    useMemo,
    useReducer,
    useRef,
    useState,
  } from 'react'

  export function useHooks<Props = {}>(func: (props: Props) => void): FC<Props>
}
