package docker

// Error codes for docker client responses
var (
	CodeBlobUnknown         = "BLOB_UNKNOWN"
	CodeBlobUploadInvalid   = "BLOB_UPLOAD_INVALID"
	CodeBlobUploadUnknown   = "BLOB_UPLOAD_UNKNOWN"
	CodeDigestInvalid       = "DIGEST_INVALID"
	CodeManifestBlobUnknown = "MANIFEST_BLOB_UNKNOWN"
	CodeManifestInvalid     = "MANIFEST_INVALID"
	CodeManifestUnknown     = "MANIFEST_UNKNOWN"
	CodeManifestUnverified  = "MANIFEST_UNVERIFIED"
	CodeNameInvalid         = "NAME_INVALID"
	CodeNameUnknown         = "NAME_UNKNOWN"
	CodeSizeInvalid         = "SIZE_INVALID"
	CodeTagInvalid          = "TAG_INVALID"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeDenied              = "DENIED"
	CodeUnsupporteD         = "UNSUPPORTED"
)
