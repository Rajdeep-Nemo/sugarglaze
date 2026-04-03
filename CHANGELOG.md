# Changelog

## [0.2.0] - 1st April 2026
- **Strictly-Typed Variable System:** Fully implemented let and const with mandatory initialization rules.

- **Memory-Safe Bounds Checking:** Evaluator now strictly catches numeric overflow/underflow for all integer and float types instead of silently wrapping.

- **Smart Assignments:** Added strict type enforcement for standard assignment (=) and container-driven implicit casting for compound assignments (+=, -=, etc.).

- **Context-Aware Errors:** Upgraded error reporting to provide exact variable names, expected types, and precise boundary-crossing details.

## [0.1.0] - 27th March 2026
- Initial implementation