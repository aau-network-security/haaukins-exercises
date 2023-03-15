package testdata

import "github.com/aau-network-security/haaukins-exercises/proto"

var (
	InSuccessJobspaceOnly = []*proto.Exercise{
		{
			Name:        "Test Challenge",
			Tag:         "test_challenge",
			Secret:      true,
			Category:    "Privacy Universe",
			Description: "some irrelevant text",
			Platforms:   []string{"jobspace"},
			PrivacyEnv:  "test",
			Hosts: []*proto.Host{
				{
					Flags: []*proto.Flags{
						{
							Name:        "Some Jobspace challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Privacy Universe",
						},
					},
				},
			},
		},
	}
	OutSuccessJobspaceOnly = []*proto.Exercise{
		{
			Name: "Privacy Universe",
			Tag:  "pu",
			Hosts: []*proto.Host{
				{
					Image: "registry.gitlab.com/haaukins/privacy-universe/jobspace",
					Dns: []*proto.DNSRecord{
						{
							Type: "A",
							Name: "jobspace.hkn",
						},
					},
				},
				{
					Image: "registry.gitlab.com/haaukins/privacy-universe/privacy-api",
					Dns: []*proto.DNSRecord{
						{
							Type: "A",
							Name: "api.jobspace.hkn",
						},
					},
					Environment: []*proto.EnvVariable{
						{
							Name:  "CHALLS",
							Value: "test",
						},
					},
					Flags: []*proto.Flags{
						{
							Name:        "Some Jobspace challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Privacy Universe",
						},
					},
				},
			},
		},
	}
	InSuccessNoPU = []*proto.Exercise{
		{
			Name:        "Test Challenge",
			Tag:         "test_challenge",
			Secret:      true,
			Category:    "Web Exploitation",
			Description: "some irrelevant text",
			Hosts: []*proto.Host{
				{
					Image: "registry.gitlab.com/haaukins/web/test",
					Flags: []*proto.Flags{
						{
							Name:        "Some web challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Web Exploitation",
						},
					},
				},
			},
		},
	}
	OutSuccessNoPU = []*proto.Exercise{
		{
			Name:        "Test Challenge",
			Tag:         "test_challenge",
			Secret:      true,
			Category:    "Web Exploitation",
			Description: "some irrelevant text",
			Hosts: []*proto.Host{
				{
					Image: "registry.gitlab.com/haaukins/web/test",
					Flags: []*proto.Flags{
						{
							Name:        "Some web challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Web Exploitation",
						},
					},
				},
			},
		},
	}
	InMixSuccess = []*proto.Exercise{
		{
			Name:        "Test Challenge1",
			Tag:         "test_challenge1",
			Secret:      true,
			Category:    "Web Exploitation",
			Description: "some irrelevant text",
			Hosts: []*proto.Host{
				{
					Image: "registry.gitlab.com/haaukins/web/test",
					Flags: []*proto.Flags{
						{
							Name:        "Some web challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Web Exploitation",
						},
					},
				},
			},
		},
		{
			Name:        "Test Challenge",
			Tag:         "test_challenge",
			Secret:      true,
			Category:    "Privacy Universe",
			Description: "some irrelevant text",
			Platforms:   []string{"jobspace"},
			PrivacyEnv:  "test",
			Hosts: []*proto.Host{
				{
					Flags: []*proto.Flags{
						{
							Name:        "Some Jobspace challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Privacy Universe",
						},
					},
				},
			},
		},
	}
	OutMixSuccess = []*proto.Exercise{
		{
			Name:        "Test Challenge1",
			Tag:         "test_challenge1",
			Secret:      true,
			Category:    "Web Exploitation",
			Description: "some irrelevant text",
			Hosts: []*proto.Host{
				{
					Image: "registry.gitlab.com/haaukins/web/test",
					Flags: []*proto.Flags{
						{
							Name:        "Some web challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Web Exploitation",
						},
					},
				},
			},
		},
		{
			Name: "Privacy Universe",
			Tag:  "pu",
			Hosts: []*proto.Host{
				{
					Image: "registry.gitlab.com/haaukins/privacy-universe/jobspace",
					Dns: []*proto.DNSRecord{
						{
							Type: "A",
							Name: "jobspace.hkn",
						},
					},
				},
				{
					Image: "registry.gitlab.com/haaukins/privacy-universe/privacy-api",
					Dns: []*proto.DNSRecord{
						{
							Type: "A",
							Name: "api.jobspace.hkn",
						},
					},
					Environment: []*proto.EnvVariable{
						{
							Name:  "CHALLS",
							Value: "test",
						},
					},
					Flags: []*proto.Flags{
						{
							Name:        "Some Jobspace challenge",
							Tag:         "test_challenge-0",
							Static:      "HKN{test}",
							Points:      10,
							Description: "some description",
							Category:    "Privacy Universe",
						},
					},
				},
			},
		},
	}
)
