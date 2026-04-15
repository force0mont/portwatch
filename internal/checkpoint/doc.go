// Package checkpoint provides persistence for the last-known set of open
// ports observed by portwatch.
//
// On startup the daemon loads the checkpoint file so it can distinguish
// ports that existed before it launched from genuinely new listeners that
// appeared while it was running.  After each successful scan cycle the
// watcher calls Save to atomically replace the on-disk state.
//
// The file is written via a rename-into-place strategy to avoid partial
// writes corrupting the stored snapshot.
package checkpoint
