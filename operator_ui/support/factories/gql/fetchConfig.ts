// buildConfigItem builds the config item fields.
export function buildConfigItem(
  overrides?: Partial<Config_ItemsFields>,
): Config_ItemsFields {
  return {
    __typename: 'ConfigItem',
    key: 'BAND',
    value: 'Major Lazer',
    ...overrides,
  }
}
