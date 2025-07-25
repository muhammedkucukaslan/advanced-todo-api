# Using `*validator.Validate` Directly in Handler Layer

- Status: accepted
- Date: 2025-05-18

## Context and Problem Statement

To ensure the correctness of HTTP requests, we need to use a validation library. We chose the github.com/go-playground/validator package. However, within our project architecture, a question arises: should we define an abstract interface named Validator in the handler layer, or directly use the *validator.Validate instance?

## Decision Drivers

- Should be easy to use
- Should not make the project more complex
- Should minimize code overhead
- Should not break the logic behind abstraction unnecessarily

## Considered Options

- `Define an abstract interface named Validator
- Use the *validator.Validate instance directly in the Validator field

## Decision Outcome

The chosen decision is to use the `*validator.Validate` instance directly, because defining an interface named Validator and creating a separate file for it in every handler struct introduces unnecessary overhead. Furthermore, since we plan to consistently use the same validation library, defining an abstract interface for it does not offer significant architectural benefits.

### Positive Consequences

- Reduced code overhead
- Cleaner and more understandable file structure

### Negative Consequences

- Trade-off in architecture.
- Testing became more difficult.
