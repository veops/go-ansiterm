# Changelog

## v0.0.6

- Add idempotent Stream.Close and stop/join parser workers on close or reattach.
- Safely cancel channel operations during concurrent shutdown and partial escape sequences.
- Fix readline deletion and erasure at nonzero columns (#2).
- Fix one-based absolute cursor positioning, including scrolling margins (#3).
- Preserve interior blanks and split UTF-8 in terminal command output.
- Accept empty OSC titles and retain their full text.
- Reduce screen rendering, scrolling and text translation allocations.
- Document parser lifecycle and command policy responsibilities (#1).
