# Berth presence

Provision each Berth device by adding its machine ID to
`XBGW_BERTH_MACHINE_IDS` in the gateway environment. The ID is the EEPROM's
12-digit decimal serial formatted as 16 lowercase hex digits. Multiple IDs are
comma-separated. Keep the operator environment private; never put EEPROM data
or credential material in source control. Future Berth onboarding must maintain
this device attribution alongside device provisioning; client discovery is not
automatic in the native Xbox protocol.

The gateway obtains the machine identity from the sealed service ticket's
`SN.<serial>@<domain>` principal after verifying both ticket and authenticator.
The normal compound machine/user exchange preserves this principal. It forwards
`client=berth` over the existing mTLS feature stream. xb-social publishes that
client in the existing `xbox` namespace, including when a core lease is renewed
or recreated. A new physical-console connection clears the emulator label.

Unlisted devices and user-only tickets keep the Xbox label. This is a label for
a provisioned identity, not emulator attestation: copying the same EEPROM to
another client copies its attribution too. Do not share a console identity
between emulator and hardware. No user/account-name or source-IP inference is
used. The shared Xbox friends, invitations and title identifiers stay intact.

Rollout: update xb-social first, then the gateway and its environment. Only those
two services need restarting; no core or Switch restart is required. Existing
connections need a new gateway handshake and Alive to publish the label; if the
client does not reconnect automatically, sign out and back in. Switch and website
already display the client label as “Berth”. Roll back by restoring the previous
two images and environment. No database migration is involved.
