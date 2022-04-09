package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnUvBg6ERqyFpWedZ8F+H
wJ2UiT7M9b6wdMPLI9Ikmrt/aYm+vUc2vqaYPrzKTupdEqGzoJI7e2g2JDJgZG/K
wiz1JD4iGifC64zw40wl5SgjP5+MqKY9C9jas7pq6lfWwvQ2a1XJURkuLg4NF0Vd
ybVqEkQc7y6WyZgKx4rjkH0LefM8k401HnVcowmmMf02cQ86bb04U2DqMGUzh5AZ
SrgFbeuAbERmi2mHzi8lJ5Okd7VFA2kA4Den2JqALyugaYvNr/oV107EgkRjMJZt
Y/F+kxiWcsLfqeXm5gcL1/78gHpC3drT43NYNssFhmnJaL9jHpChMRZjoyDa7i+4
iQIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}

	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}

	// 表格驱动测试
	cases := []struct {
		name    string
		token   string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name:  "valid_token",
			token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI0OTNhZTk2OTBmNTAzOTJmMjRjZmFkIn0.Php6cc_AFplJKyLw-RdNUAe8wcamYU11TgYpBgzasv5f52MHNYvFvLNHcUCMyeUNSwx52jMsR-_NDYfzXgHuMsmcUmttRLCOKHPA0t-wn821ump16LjQg8G6MV9suGsbCJY363Z7gDAGBYk4lVe6yT9zxFwscvR6bEVvLRyEXk-NbcuEcARIEcEg-KV00uGPQ4_UpH9Wgjc6MGmAsd6oz-mDoR2n1Wz41CRKfpdQM34NSijIUs5ASonJQf92B0fwO42fXLzQYuF7fg6r6NLMhLtMCGiJ5aTvjELfFz8Myu_vvbwldddBXPUTfZILmqo3cfaz42jPgl_QAqzYacW45w",
			now:   time.Unix(1516239122, 0),
			want:  "62493ae9690f50392f24cfad",
		},
		{
			name:    "token_expired",
			token:   "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI0OTNhZTk2OTBmNTAzOTJmMjRjZmFkIn0.Php6cc_AFplJKyLw-RdNUAe8wcamYU11TgYpBgzasv5f52MHNYvFvLNHcUCMyeUNSwx52jMsR-_NDYfzXgHuMsmcUmttRLCOKHPA0t-wn821ump16LjQg8G6MV9suGsbCJY363Z7gDAGBYk4lVe6yT9zxFwscvR6bEVvLRyEXk-NbcuEcARIEcEg-KV00uGPQ4_UpH9Wgjc6MGmAsd6oz-mDoR2n1Wz41CRKfpdQM34NSijIUs5ASonJQf92B0fwO42fXLzQYuF7fg6r6NLMhLtMCGiJ5aTvjELfFz8Myu_vvbwldddBXPUTfZILmqo3cfaz42jPgl_QAqzYacW45w",
			now:     time.Unix(1517239122, 0),
			wantErr: true,
		},
		{
			name:    "bad_token",
			token:   "bad_token",
			now:     time.Unix(1517239122, 0),
			wantErr: true,
		},
		{ // 伪造的token,AccountID尾部的ad改成a4，伪造的用户不知道privateKey，用ad生成的key去伪造
			name:    "wrong_signture",
			token:   "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI0OTNhZTk2OTBmNTAzOTJmMjRjZmE0In0.Php6cc_AFplJKyLw-RdNUAe8wcamYU11TgYpBgzasv5f52MHNYvFvLNHcUCMyeUNSwx52jMsR-_NDYfzXgHuMsmcUmttRLCOKHPA0t-wn821ump16LjQg8G6MV9suGsbCJY363Z7gDAGBYk4lVe6yT9zxFwscvR6bEVvLRyEXk-NbcuEcARIEcEg-KV00uGPQ4_UpH9Wgjc6MGmAsd6oz-mDoR2n1Wz41CRKfpdQM34NSijIUs5ASonJQf92B0fwO42fXLzQYuF7fg6r6NLMhLtMCGiJ5aTvjELfFz8Myu_vvbwldddBXPUTfZILmqo3cfaz42jPgl_QAqzYacW45w",
			now:     time.Unix(1516239122, 0),
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := v.Verify(c.token)

			if !c.wantErr && err != nil {
				t.Errorf("verification failed: %v", err)
			}

			if c.wantErr && err == nil {
				t.Errorf("want error; got no error")
			}

			if accountID != c.want {
				t.Errorf("wrong account id. want: %q; got: %q", c.want, accountID)
			}

		})
	}

	// token := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI0OTNhZTk2OTBmNTAzOTJmMjRjZmFkIn0.Php6cc_AFplJKyLw-RdNUAe8wcamYU11TgYpBgzasv5f52MHNYvFvLNHcUCMyeUNSwx52jMsR-_NDYfzXgHuMsmcUmttRLCOKHPA0t-wn821ump16LjQg8G6MV9suGsbCJY363Z7gDAGBYk4lVe6yT9zxFwscvR6bEVvLRyEXk-NbcuEcARIEcEg-KV00uGPQ4_UpH9Wgjc6MGmAsd6oz-mDoR2n1Wz41CRKfpdQM34NSijIUs5ASonJQf92B0fwO42fXLzQYuF7fg6r6NLMhLtMCGiJ5aTvjELfFz8Myu_vvbwldddBXPUTfZILmqo3cfaz42jPgl_QAqzYacW45w"
	// jwt.TimeFunc = func() time.Time {
	// 	return time.Unix(1516239122, 0)
	// }
	// accountID, err := v.Verify(token)

	// if err != nil {
	// 	t.Errorf("verification failed: %v", err)
	// }

	// want := "62493ae9690f50392f24cfad"
	// if accountID != want {
	// 	t.Errorf("wrong account id. want: %q; got: %q", want, accountID)
	// }

}
