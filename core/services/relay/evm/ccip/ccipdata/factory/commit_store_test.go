package factory

//
//func TestCommitStore(t *testing.T) {
//	ctx := t.Context()
//	for _, versionStr := range []string{ccipdata.V1_2_0} {
//		lggr := logger.Test(t)
//		addr := cciptypes.Address(utils.RandomAddress().String())
//		lp := lpmocks.NewLogPoller(t)
//
//		feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)
//
//		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
//		versionFinder := factory2.newMockVersionFinder(ccipconfig.CommitStore, *semver.MustParse(versionStr), nil)
//		_, err := NewCommitStoreReader(ctx, lggr, versionFinder, addr, nil, lp, feeEstimatorConfig)
//		assert.NoError(t, err)
//
//		expFilterName := logpoller.FilterName(v1_2_0.ExecReportAccepts, addr)
//		lp.On("UnregisterFilter", mock.Anything, expFilterName).Return(nil)
//		err = CloseCommitStoreReader(ctx, lggr, versionFinder, addr, nil, lp, feeEstimatorConfig)
//		assert.NoError(t, err)
//	}
//}
