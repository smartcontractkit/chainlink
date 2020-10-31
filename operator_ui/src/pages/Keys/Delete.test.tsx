import React from 'react'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import { Delete } from './Delete'

describe('pages/Keys/Delete', () => {
  it('open modal and confirm delete', async () => {
    const expectedKeyId = 'KeyId'
    const expectedKeyValue = 'keyValue'
    const expectedOnDelete = jest.fn()

    const wrapper = mountWithProviders(
      <Delete
        onDelete={expectedOnDelete}
        keyId={expectedKeyId}
        keyValue={expectedKeyValue}
      />,
    )

    wrapper.find('[data-testid="keys-delete-dialog"]').first().simulate('click')

    expect(wrapper.text()).toContain(expectedKeyValue)

    wrapper
      .find('[data-testid="keys-delete-confirm"]')
      .first()
      .simulate('click')

    expect(expectedOnDelete).toBeCalledWith(expectedKeyId)
  })
})
