package deployments_test

// func TestGetVersion(t *testing.T) {
// 	inputs := models.Inputs{
// 		Source: models.Source{
// 			ServerURL: "",
// 			ProjectID: "",
// 			Auth: models.SourceAuthentication{
// 				Token: "",
// 			},
// 			Deployments: &models.SourceDeployment{
// 				Environment: gitlab.Ptr("testing"),
// 				Status:      gitlab.Ptr("running"),
// 			},
// 			Feature: "deployments",
// 		},
// 		// Version: "",
// 	}
//
// 	expectUser := gitlab.ProjectUser{
// 		ID:       1,
// 		Username: "user",
// 	}
// 	expectEnvironment := gitlab.Environment{
// 		ID:   1,
// 		Name: "testing",
// 	}
// 	expectedResp := &gitlab.Response{}
// 	pid := 1
//
// 	client := gitlabtesting.NewTestClient(t)
// 	opts := &gitlab.ListProjectDeploymentsOptions{
// 		Sort:        gitlab.Ptr("desc"),
// 		Environment: inputs.Source.Deployments.Environment,
// 		Status:      inputs.Source.Deployments.Status,
// 	}
//
// 	expectDeployments := []*gitlab.Deployment{
// 		{
// 			ID:          3,
// 			Status:      *inputs.Source.Deployments.Status,
// 			Environment: &expectEnvironment,
// 			Ref:         "main",
// 			SHA:         "aafea2f49feazduahziudhazidhaiuzdhajzhdakl",
// 			IID:         3,
// 			User:        &expectUser,
// 		},
// 		{
// 			ID:          2,
// 			Status:      *inputs.Source.Deployments.Status,
// 			Environment: &expectEnvironment,
// 			Ref:         "main",
// 			SHA:         "aafea2f49feb3f256906e88f7d614362c99b3856",
// 			IID:         2,
// 			User:        &expectUser,
// 		},
// 		{
// 			ID:          1,
// 			Status:      *inputs.Source.Deployments.Status,
// 			Environment: &expectEnvironment,
// 			Ref:         "main",
// 			SHA:         "4ae9025b3a99e3cdd7c09dc4387895658df20a3a",
// 			IID:         1,
// 			User:        &expectUser,
// 		},
// 	}
//
// 	client.MockDeployments.
// 		EXPECT().
// 		ListProjectDeployments(pid, opts).
// 		Return(expectDeployments, expectedResp, nil)
//
// 	deploys, _, err := client.Deployments.ListProjectDeployments(
// 		pid, opts,
// 	)
// 	assert.NoError(t, err, "list deployments from mock should not fail")
//
// 	t.Run("CheckFirstVersionOk", func(t *testing.T) {
// 		version, err := deployments.GetVersion(deploys, &inputs)
// 		assert.NoError(t, err, "Get Version should not fail")
// 		assert.NotNil(t, version)
//
// 		versionDeploy, ok := version.(*gitlab.Deployment)
// 		assert.True(t, ok, "type conversion should work")
// 		assert.Equal(t, expectDeployments[0].ID, versionDeploy.ID, "id should match the one of the last version")
// 	})
//
// 	t.Run("CheckUpdateOk", func(t *testing.T) {
// 		inputs.Version = expectDeployments[1]
// 		version, err := deployments.GetVersion(deploys, &inputs)
// 		assert.NoError(t, err, "Get Version should not fail")
// 		assert.NotNil(t, version)
//
// 		versionDeploy, ok := version.([]*gitlab.Deployment)
// 		assert.True(t, ok, "type conversion should work")
// 		assert.ElementsMatch(t, expectDeployments[:1], versionDeploy, "element should match from the 2nd item of the list")
// 	})
// }
