resource "databasus_users_settings" "mySettings" {
  allow_external_registrations        = false
  allow_member_invitations            = false
  member_allowed_to_create_workspaces = false
}
