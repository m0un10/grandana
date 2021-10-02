resource "google_storage_bucket" "tfstate" {
  name          = "${var.env}-${var.state_name}-state-file"
  uniform_bucket_level_access = true
  location = "US"
  encryption {
     default_kms_key_name = "${google_kms_crypto_key.gcs.self_link}"
  }
  versioning {
    enabled = true
  }
  depends_on = [google_project_iam_member.grant-google-storage-service-encrypt-decrypt]
}

resource "google_kms_crypto_key" "gcs" {
  name            = var.gcp_key_name
  key_ring        = "${google_kms_key_ring.gcs.self_link}"
  rotation_period = "86401s"
}

resource "google_kms_key_ring" "gcs" {
  name     = var.gcp_key_ring
  location = "us"
}

# grant KMS Encrypter Decryptor permissions to the Google Storage Service
resource "google_project_iam_member" "grant-google-storage-service-encrypt-decrypt" {
  role   = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.current.number}@gs-project-accounts.iam.gserviceaccount.com"
}