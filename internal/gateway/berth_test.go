// SPDX-License-Identifier: AGPL-3.0-only
package gateway

import (
	"github.com/openpak/xb-common/krb"
	"github.com/openpak/xb-gateway/sgwire"
	"testing"
	"time"
)

func TestBerthUsesAuthenticatedMachinePrincipal(t *testing.T) {
	ring := sgRing(t, "a1b2c3d4e5f60718293a4b5c6d7e8f90")
	for _, tc := range []struct {
		name string
		want string
	}{
		{"SN.123456789012@xbox.com", "0000001cbe991a14"},
		{"Player@xbox.com", ""}, {"SN.123@xbox.com", ""},
		{"SN.000000000000@xbox.com", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			now := time.Now()
			tk, key, _ := issueTicket(t, ring, now, advertise, time.Hour)
			part, err := krb.OpenTicket(ring, krb.UsageSGTicket, tk)
			if err != nil {
				t.Fatal(err)
			}
			part.CName = krb.PrincipalName{Type: krb.NameTypeEnterprise, Parts: []string{tc.name}}
			kv, raw := ring.Active()
			sealKey, _ := krb.ToKey(raw)
			tk.Enc, err = krb.Seal(sealKey, krb.UsageSGTicket, &kv, part.Marshal())
			if err != nil {
				t.Fatal(err)
			}
			packet := realInit(t, tk, key, part.CName, now)
			init, err := sgwire.ParseInit(packet)
			if err != nil {
				t.Fatal(err)
			}
			id, err := (KrbVerifier{SGKeys: ring}).VerifyAPReq(init.APReq, now)
			if err != nil || id.MachineID != tc.want {
				t.Fatalf("identity=%+v error=%v", id, err)
			}
			// Changing only the client's authenticator cannot impersonate a machine.
			fakeName := krb.PrincipalName{Type: krb.NameTypeEnterprise, Parts: []string{"SN.999999999999@xbox.com"}}
			forged, _ := sgwire.ParseInit(realInit(t, tk, key, fakeName, now))
			if _, err := (KrbVerifier{SGKeys: ring}).VerifyAPReq(forged.APReq, now); err == nil {
				t.Fatal("unsealed machine attribution accepted")
			}
		})
	}
}
func TestBerthDeviceConfiguration(t *testing.T) {
	for _, tc := range []struct {
		value string
		valid bool
	}{
		{"0000001cbe991a14", true}, {" 0000001cbe991a14 , 0000000000000001 ", true},
		{"0009000000000001", false}, {"0000000000000000", false}, {"garbage", false}, {"0000001cbe991a14,", false},
	} {
		env := map[string]string{"XBGW_ADVERTISE_IP": "203.0.113.5", "XBGW_BERTH_MACHINE_IDS": tc.value}
		cfg, err := ConfigFromEnv(func(k string) (string, bool) { v, ok := env[k]; return v, ok })
		if (err == nil) != tc.valid {
			t.Fatalf("valid=%t error=%v", tc.valid, err)
		}
		if err == nil {
			s := &Server{cfg: cfg}
			if s.clientForDevice("0000001cbe991a14") != "berth" || s.clientForDevice("") != "" || s.clientForDevice("0000000000000002") != "" {
				t.Fatal("device-specific labeling failed")
			}
		}
	}
}

func TestBerthDeviceLabelReachesFeatureContext(t *testing.T) {
	for _, device := range []string{"0000001cbe991a14", "0000000000000002", ""} {
		h := newHarness(t, func(c *Config) { c.BerthMachineIDs = map[string]bool{"0000001cbe991a14": true} })
		c := h.connect(func(id *Identity) { id.MachineID = device })
		c.tcp(sgwire.TCPHeader{Src: 4000, Dst: 107, Seq: 1000, FlagsAndHeader: 6<<12 | sgwire.TCPSyn, Window: 8192}, []byte{2, 4, 0x05, 0x18})
		if _, _, ok := c.readTCP(time.Second); !ok {
			t.Fatal("no feature connection")
		}
		h.backend.mu.Lock()
		if len(h.backend.ctxs) != 1 {
			h.backend.mu.Unlock()
			t.Fatal("feature not reached")
		}
		label := h.backend.ctxs[0].Client
		h.backend.mu.Unlock()
		want := ""
		if device == "0000001cbe991a14" {
			want = "berth"
		}
		if label != want {
			t.Fatalf("client=%q want=%q", label, want)
		}
	}
}
