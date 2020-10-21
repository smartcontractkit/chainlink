import React from 'react'

export const useLoadingPlaceholder = (
  isLoading = false,
): {
  isLoading: boolean
  LoadingPlaceholder: React.FC
} => {
  const LoadingPlaceholder: React.FC = isLoading
    ? () => <div>Loading...</div>
    : () => null

  return {
    isLoading,
    LoadingPlaceholder,
  }
}
