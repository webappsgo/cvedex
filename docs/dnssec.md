# DNSSEC

DNSSEC key generation, signing, signature refresh, and DS publication are not implemented in the current build. The `/api/v1/ds` route therefore returns HTTP 503. Do not publish a parent-zone delegation that assumes this server signs its zone.

DNSSEC support is part of the project specification and remains required implementation work.
