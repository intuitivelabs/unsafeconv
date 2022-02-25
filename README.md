# unsafeconv

This module contains fast conversion functions between string and byte slices,
 with no copy or allocation.
These functions should be used only if the byte slice content is never modified
(treated as read-only).
