package offchainreporting

func (c *ConfigOverriderImpl) ExportedUpdateFlagsStatus() error {
	return c.updateFlagsStatus()
}
