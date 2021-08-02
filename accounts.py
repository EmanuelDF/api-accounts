import requests, time, uuid, base64 

from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives import serialization, hashes

# test key pair
# public key is stored inside Form3 system and identified by the keyid
test_public_key_id="75a8ba12-fff2-4a52-ad8a-e8b34c5ccec8"
test_private_key_filename="test_private_key.pem"

# Replace these variables with your own data! ###
organisation_id = 'YOUR ORGANISATION ID HERE'
bank_id = 'YOUR UK SORTCODE HERE'

url = "/v1/http-signatures/debug"
host = 'api.staging-form3.tech'
base_url = 'https://%s' % host

# Generate IDs
payment_id = uuid.uuid4()

# take current time
time_now = time.gmtime()
# time formatted for payload
request_date = time.strftime("%Y-%m-%d", time_now)
# time formatted for signature
signature_date = time.strftime("%a, %d %b %Y %T %Z", time_now)

# full address to our endpoint
signature_debug_endpoint = "%s%s" % (base_url, url)

payment_payload = """
{
    "data": {
        "type": "payments",
        "id": "%s",
        "version": 0,
        "organisation_id": "%s",
        "attributes": {
            "amount": "200.00",
            "beneficiary_party": {
                "account_name": "Ada Lovelace",
                "account_number": "71268996",
                "account_number_code": "BBAN",
                "account_with": {
                    "bank_id": "400302",
                    "bank_id_code": "GBDSC"
                }
            },
            "currency": "GBP",
            "debtor_party": {
                "account_name": "Isaac Newton",
                "account_number": "87654321",
                "account_number_code": "BBAN",
                "account_with": {
                    "bank_id": "%s",
                    "bank_id_code": "GBDSC"
                }
            },
            "processing_date": "%s",
            "reference": "Something",
            "payment_scheme": "FPS",
            "scheme_payment_sub_type": "TelephoneBanking",
            "scheme_payment_type": "ImmediatePayment"
        }
    }
}
""" % (payment_id, organisation_id, bank_id, request_date)

with open(test_private_key_filename, "rb") as key_file:
    test_private_key = serialization.load_pem_private_key(
        key_file.read(),
        # Private key in this example is not protected by a passphrase
        # If you protect your private key with a passphrase you will need to provide it here
        password=None,
        backend=default_backend()
    )

    # generating digest from the paylod
    digest_genererator = hashes.Hash(hashes.SHA256(), backend=default_backend())
    digest_genererator.update(payment_payload)
    body_digest = digest_genererator.finalize()

    # we cannot sent digest in binary format we need to encode it using base64
    base64_body_digest = base64.b64encode(body_digest)

    # creating signature
    # this is a list of selected headers + special header (requet-target)
    # each line has to end with `\n`
    # `digest` and `content-length` is only required when sending request with BODY
    signature = """(request-target): post %s
host: %s
date: %s
content-type: application/json
accept: application/json
digest: SHA-256=%s
content-length: %s""" % (url, host, signature_date, base64_body_digest, len(payment_payload))

    # list of headers that are part of signature
    # `headers` value and headers in signature, they need to be in sync and in order
    headers = "(request-target) host date content-type accept digest content-length"

    # signing the signature using private key
    signed_signature = test_private_key.sign(
        signature,
        padding.PKCS1v15(),
        hashes.SHA256()
    )

    # we cannot sent signature in binary form that it why we need to encode it
    base64_signed_signature=base64.b64encode(signed_signature)

    # generating authorization header
    # `rsa-sha256` is the only supported algorithm
    authorization_header = """Signature: keyId=\"%s\",algorithm=\"rsa-sha256\",headers=\"%s\", signature=\"%s\"""" % (test_public_key_id, headers, base64_signed_signature)

    base64_signature = base64.b64encode(signature)

    # base headers for the request
    payment_headers = {
        'accept': "application/json",
        'content-type': "application/json",
        'date': signature_date,
        'authorization': authorization_header,
        'digest': base64_body_digest,
        # for debug endpoint we need to send additional header signature-debug
        # this header should never be used on other endpoints or production
        # it is used to generate a diff of what was sent and what was expected by the server
        # it is a base64 encoded raw signature (not signed)
        'signature-debug': base64_signature
    }

    # sending request to debug endpoint
    response = requests.request("POST", signature_debug_endpoint, data=payment_payload, headers=payment_headers)

    # response should have message `hi` and body containing the request
    print(response.text)
