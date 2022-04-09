package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwEU2vILMbXtUEmj/BSRjmxvJhe9TIw5OmYevXKRCm8EFmFdr
TQRpOn+dL/259qxpomf+ceSH6DN5QB4JWoWcb+cyo+Uqk76gzG+mUQqI1a6y5xAG
80JArWgRr9o1WRb5Es5f9LlmYMqtXI+v7FN73mlglWfgOsSnD92Tzc+U3KghR7G6
FbTw2NiZXd4pCXflcQtXQUGvZIgW1tHP9yw3TT2TtAP5rH7pZmOI4vwLgfL9IlHe
oJhdZPSyTbt5FQFJz0L4yVZ1EQCyibdnAicUTtaznIvx6yQ7ZenIx6lg/Hri3Jb8
sSzRnSoyiSOdcHlDHEwpWDDDJko5S+SaZiMxzwIDAQABAoIBACtj81RbmFLk0DZP
Z637Zlcs0C/fsL1qjyZqzpJzp2yLBP46TEbXRgAjqI+aDQu0ISq7mVqOnnUymJx1
qtq46HMAlBcUsqUhEhzQ3ZHQdMz0Xf9zSH5BZ2M65zLuZbezTYaG+HS8GnShiLoo
2yTUOlIZKejNxna8xo74DFEYFZouGvZ6aS8QLGD5woD5m0D5WHkDaejgsxAaWFq2
YH/qho0ycslmdDmORYNx0rSzaHE1gJ9ObGZxdNTqGS16JxLgV5TCovlp09b2AvpX
lVmhRZ19gUmLBM/bV4vYwwx6h93FPGLIQWnyAZS++GSF2Q2OJYNU0MkF4a6om1H5
webVSUECgYEA6xhuXlfQYRfw8nqsPNOIOIQBhd/iCuHhq1eqRNhDz2XmtrGckjVm
4FGAHq2dex8YzA1yNRwtW0u4gIzerVFe9hEDv6RTJsBA2lArPI4UFlAvchTidb0H
/ce5ffaDfaPKeGV+LvK4OMyMk/5jQ+IWUBz4dSXSbktwS3N49BFSjZECgYEA0V3w
Yok3IpVIqRcevYkc5WqQrHAcPg/VHJWBSXICA9zMUopCDs7yHDrx6s4i93fDNVNX
Tpht/3YQgT5dQPZ/fla8RJouUHuTDM2rDOc8ygUbhjp+RpOJ0SOiuUvOinoTjyAX
PhNrBBVA7WA8f0zVe0wjeqE9BEobxdn+iTUVmV8CgYEAwBTECF0jKlFPUi6cj04d
nF9uhx03KMEJit8jhQBV1xxg9qADVwo+tcShM2+Sncf4kK/jwjT1cgRaCN4QCffT
6KRaNbhjH+QwmK6rxNwo2VpGiaU6qgv6fYUzrxE8ueibQudZw64Yin/F1B0iqZIN
vQMOlViDqA4G//6msnN8yaECgYBwx8Cg47duiqoMY+bsRHrrXh82tGGSUOcEschm
LuLE/+CUqeNxpKqo6Fu+l3IgikP+armCHfxK+2dip5yGTQJapRBfiApg0mBhKjz1
A28sh9nO0Z2KGRnJLgAO/rXwxFfa5nd+uekQ1v4VoJyWGmZ5N4d5HHgI4n7Zcld5
w72x0wKBgQDMytwsp89LXnO/yL+n28R2EqmaZOeYUVoG4F29m4Li5ucPacOhGcqQ
paWDi1YliYmB+nseVrVDSgsshir4BqP2Fd7nwKN4hg1s5YDy9nkiFaW3lICSk9aN
cjCigazvtWLFZ6qIf7Hv3a4g8ZxvrScO/fY2JIuSmflJ17Y3pIvpFw==
-----END RSA PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key: %v", err)
	}

	g := NewJWTTokenGen("coolcar/auth", key)
	g.nowFunc = func() time.Time {
		// 固定下来，每次的测试才有意义
		return time.Unix(1516239022, 0)
	}

	token, err := g.GenerateToken("62493ae9690f50392f24cfad", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token: %v", err)
	}

	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI0OTNhZTk2OTBmNTAzOTJmMjRjZmFkIn0.O7qbK2TNMchnJVG6knnk-GGr-zMYm50X8q_EaQ-eRivtOU2-hyDd5LWZdwOp-VOD4TKCwaEBOH4gvJTN72W2eeVMn_NIcZ7N2gd9QJd59xL6CcYTbloPoYt-K04HB9RuuJUygqpJQib80kyEogRhxuCdAjXTjZypRAXGa-T6Bddx2c2gPbmmX-h_FWa6Gz_pLJcwtauS1TB25EnYl7HuA2YQMXTWqQdGojhlnzAIV1weG1aQCkshjpWzta0fzuhRT5GF5YSDM5So8gdjoiWUYllTFqyM0BIdOUoirA5-7gqKuKNbdCg9cb0YsQLIo7eJEJJoInj1yslRMVTHsGPAXQ"
	if token != want {
		t.Errorf("wrong token generated. want: %q; got: %q", want, token)
	}

}
