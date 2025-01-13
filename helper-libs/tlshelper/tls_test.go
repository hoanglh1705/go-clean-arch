package tlshelper

import (
	"fmt"
	"testing"
)

func TestExtractPublicKeyFromCertificateSuccess(t *testing.T) {
	tlsConfig, err := NewClientTLSConfigFromBase64(
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUV6VENDQXJXZ0F3SUJBZ0lVQ1hvbFZsN2NlLy81bXIwRy82blhHVUlCK2tNd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1h6RUxNQWtHQTFVRUJoTUNWazR4RERBS0JnTlZCQWdNQTBoRFRURVBNQTBHQTFVRUJ3d0dWR2gxUkhWagpNUTB3Q3dZRFZRUUtEQVJEYVcxaU1SQXdEZ1lEVlFRTERBZERhVzFpTFVsVU1SQXdEZ1lEVlFRRERBZERhVzFpCklFTkJNQjRYRFRJME1EUXlOVEE0TURZd01Gb1hEVEkzTURFeU1EQTRNRFl3TUZvd2FURUxNQWtHQTFVRUJoTUMKVms0eEREQUtCZ05WQkFnTUEwaERUVEVNTUFvR0ExVUVCd3dEU0VOTk1RMHdDd1lEVlFRS0RBUkRhVzFpTVJBdwpEZ1lEVlFRTERBZERhVzFpTFVsVU1SMHdHd1lEVlFRRERCUXFMbU5wYldKa1pYWXRaR1YyTFdOc2FXVnVkRENDCkFTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBSitmaXIzemlRbWcvU1RwWjFJRENiR2wKN0NFelNHYVo4MnJtVWtoSmtGbGk1bTRRbDBaNmllY1lBSDhkbStvakFNT2RxWDMyUGxGUlNxQlBPUVNwRW91QwpBSTBrYXBiR0UzeTFLTDhWQSs4ejN1bW1vWHNiYW5yKy9Ga1RqVEdmZzJqNVdyRWpiUSs1VUFIbS9ib1I4K0FMCkRsQktEaEdPcDJOSGQ1OTV4a0Y4RU04a01FbW0yTjNORW9lNkV4cXdJdWZ3djBYVVFEVlk2dExENG1wTWZWcjYKeWsvNnFvenBRb0pZUU4ydkZJN1NJTFFlWEJoclFFZjVyK0o5VklyK1JQcDFKaHhFdlEvRU9CY0l4bnVkcmxqQwp3elBEV003dnBvTW5NNTZPZVY2eFQxMWpNMnZJVkRPMHZGNW5UTCszNmNZMjZLRHltVmxJTXdFbE1RekEwVzhDCkF3RUFBYU4zTUhVd013WURWUjBSQkN3d0tvSVNZMmx0WW1SbGRpMWtaWFl0WTJ4cFpXNTBnaFFxTG1OcGJXSmsKWlhZdFpHVjJMV05zYVdWdWREQWRCZ05WSFE0RUZnUVVmdE1pNW45SU1RSm44R29qbUVIa3N0U0RVNU13SHdZRApWUjBqQkJnd0ZvQVUxWTlQUlJoWUFldHh1T1BGcGhHQ3BtN1JjY013RFFZSktvWklodmNOQVFFTEJRQURnZ0lCCkFBV2dWRHVFVFhmMy9SakpHczNBMWtQWVNYdEt3Qm44K09nM2xFS3RTUVJNdHYvS3J1VGF0RUtBZTY1Z295dnEKTUhINlRLRHRZY3V0bjMwTENGd1RhN1RoYjl3YXVkMHZhL3p2U0FXRXo4NFNYOEN3MmdxVVRNU2psRlQreHdtVAp1YWkzbU5ha0dBNDU3OUxCVTlUakFuMXdURXNsMEhzNE5ZTHd2VlhLbVRFK1FVdjFzcGppWmdpOTlhZ3JpSjY1CnhYd3dkc0UzV0pEQjlZU0xSVGJTWDFkUUNxSlVvdERnWGxqZlF1bGY5bmd6d3A0R1luOTRzeG9OcVlCN05HK3QKYzFocFk2TGpuSDRNd2E1ZkM0aGo1ZU1MZ1FFbWtWVHYwbWxUbWQxTitZams4Y0pBSWhtcS8zenloUVZ2eFRJSwp1YjJyS090WjZxcnFZU3ZxMXpyY0JTamZTV2xlYndLOEJkYThQVlliRDVCejNKV3pmOHBENXNNcXlNb2o4QkM3Ci9BRWFKR1R3OUF1aGpDV203OStqMEYyenBINDhab1hMdlZsNGhWZm4xNUErUDJJNXVveFVTZzhpTnEzL29NZlIKTjNSZzNTc3h0aHRiZlM1R3A4ZU0rWVE5V012b1owQ1dGSGM3UkIxYWQzbWNHczBzQWZ5UXRDYXhLVGdlS0lBWQpmMEF3ZlJEbW5XNU5UTWpJSnBwWXBxMFpGWUg5VGNBNUczZE1rc3Y1UW5obTFUYTh2UUEvdnZra1lFR3VYVll3CmhJSkhaVCs4OFlwY2ZPNVRIbFJrSVNzbFJJSW1LaGx0SE42K2Z2ZnlnRlJzZjhxMXpnUjVmU2lFME45Ylk3QysKNjVONU9tMGRPSmRHUTFBWVBybWU4MTdlTVkzb25TRlFpeDkwVEM0RVlmZFkKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
		"LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2d0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktrd2dnU2xBZ0VBQW9JQkFRQ2ZuNHE5ODRrSm9QMGsKNldkU0F3bXhwZXdoTTBobW1mTnE1bEpJU1pCWll1WnVFSmRHZW9ubkdBQi9IWnZxSXdERG5hbDk5ajVSVVVxZwpUemtFcVJLTGdnQ05KR3FXeGhOOHRTaS9GUVB2TTk3cHBxRjdHMnA2L3Z4WkU0MHhuNE5vK1ZxeEkyMFB1VkFCCjV2MjZFZlBnQ3c1UVNnNFJqcWRqUjNlZmVjWkJmQkRQSkRCSnB0amR6UktIdWhNYXNDTG44TDlGMUVBMVdPclMKdytKcVRIMWErc3BQK3FxTTZVS0NXRURkcnhTTzBpQzBIbHdZYTBCSCthL2lmVlNLL2tUNmRTWWNSTDBQeERnWApDTVo3bmE1WXdzTXp3MWpPNzZhREp6T2Vqbmxlc1U5ZFl6TnJ5RlF6dEx4ZVoweS90K25HTnVpZzhwbFpTRE1CCkpURU13TkZ2QWdNQkFBRUNnZ0VBQzRvN081Ulp1ZXA3b3FtRklMYTdncTVlTGVCSlFiR3JtWFRoU2Z5WGhQN1QKYmUyaElpVkZ4d0ZETisxcUVqbEptdHJSSUJ3blVUV3hVWG1vdzU5OWlieHlVY0hxT3RCREpHYnNkVFFOdnNOZgpRRUVkdDRxNTNmNkZPK05mOTlCeWJhcHBWaGtiajJGMFdVN2IxUkhyTWExYThZOEpDVmZvM2hLU29XTHFaRjE0CmNYbE82MlVQUmRPTG41Si9IQWtHWGpQdXdKV2NhbVlKcHFLcko5UVZkUWwrTHhuWGdpRlZKN1J6aFJVQzRGemgKZHpROHkrVlNlRFhmRTQ2bDBhZGsyUHRlK3llMlIveXdkUGpKVXEvUEVYTUk2b0k0RnU3OXNxajJkZW1SaDVYTwpJczFlcGpSQzArRE51L2h2emh3T1RwOVphZERoQU01L0F4LzZEOHVOMFFLQmdRRGZjQlNrMjRFRWhDeXArYVJnClh0bWtRY3BCTWNyV1hQdmgxejBVdnd6aWJadVZzNE40dnlvaGcyck5EOEZ6M0t2cHl6UHFTMUtYa1RhZGpJYjgKU05VanhtU0dOOGRhdnRCcitZZEJHVXdOMnM4SGlhT1RmejAycHVybEZ5TGJKT3dvOEVNSWdZRXV6Y2FuUlRyLwpQZk9zZ0hzbExVMytRMGZWVFNVcU1ZTkNjd0tCZ1FDMjRyRzVwSnpGTEdtSzI4WnptTFhXaTZpbitvT3ZweXVsCjJsV0p4MkxEbXpGdDJpcjZRdGREODdYWGtqL0I1NmRBRnk2OWNBeFI5d3B1YjVhRzdOVjJWaGZDT2x5cjdFc3kKSlFETmZYMU1JcUk3OVprelhWelJQbElUSjhpVStsK1dJVXJoeDZZVHhzWWV5RG1iU0VqcTVHTkNKV1BWQmtCWgpTd1NiL3JTcUZRS0JnUUNWMjROTkNwZWVvdE9kOFB4cFVsUmdrV3VJakUvQnREMlB2QitRY1k2L3NzbmQvcmRYCjNjYnhFVVlwWUw3YjZZNDMwUHp4MERFRnpQUTNlTS8wRnhrMDFGUUpuUkdNOEZ1emYzbFNsUmZvVnUveDIwT2wKb25vNDFIekl4OXF0NWphcVFuS0RHdkM4cG5EdE1VYWZlRHFkWU5LM0hZcW8xUkV4bzNzZ3NIS2J0d0tCZ1FDRQpicG5sNVdiRWZRbWNUTk5pNThWZEs5cWdjUTZreHJnYnJJUGVkbXgxV3M4clRoMXJCYlhkOWYvS3I1UE50Ukx1Cm5ScnlnTTNiR0xvTUNIQUhHajdsSnlpak5DSGhPUVdtdFJia3RxZGgxMzZGVHE1MmZIZjI2VnNEbGY1d3F2RkcKeEtyMTNkM01XbGNpK1RpRjBvMUMwc2x4bjZPd0lZdTlYVTVrSzhmbGxRS0JnUUNheWZUNWNBMWk0OVBCWlpSYQpXRExLNU82RS9nZnlQMnBST0dWNjNlR1NUWjkxcndnSDBFckNsQVo5aGdVaGs2QkgyQ3RLSGhLeWM5RVlmVm5aCkpscW9aSlFuWDZPVzNuZmp6alVHMWtMLzlKNk5qczJRWFFaU1dvUGtQWGlkTlJPMk42a2NyZEFKeml2VHZhc0sKa0lTb3pQTGJzOUZRcm9UZUlGUUZyMzFiM1E9PQotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg==",
		"",
		true,
	)

	if err != nil {
		fmt.Printf("failed to create certificate, err: %v", err)
	} else {
		certificate := &tlsConfig.Certificates[0]
		cert := GetPublicKeyFromTlsCertificate(certificate)
		if cert == nil {
			fmt.Printf("failed to parse certificate, err: %v", err)
		}
	}

	if err != nil {
		t.Errorf("failed to extract public key, err: %v", err)
	}
}
