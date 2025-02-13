# Versioning

How do we use versioning in Kubernetes and not lose our minds.

## Traditional Semantic Versioning

In standard semantic versioning we use two or three numbers to signify
a version. So we have:

MAJOR.MINOR.PATCH

Where MAJOR is the major version. Different MAJOR versions are allowed to break
the API. MINOR versions are increments over a MAJOR, and must not break the API.
PATCH are updates to the MINOR - largely bug fixes.

This suggests the version is somewhat fluid, and we may be working on bugfixes
on older versions that are still in use. I like the principles here, but the idea
of maintaining mulitple versions because clients don't update to the latest seems
off to me.

## Kubernetes Semantic Versioning

Our approach is to use the following pattern:

MILESTONE.STABLE.UNSTABLE

In our approach, we drop the notion of a bugfix completely - bug fixes are made by
a fix forward in the STABLE channel. So all projects start with a MILESTONE of zero - this
is considered pre-release and should not have a stable API. a MILESTONE release is a manual
one: we decide what the requirements for the next milestone and don't push until we reach
there. We then manually signify a MILESTONE has been met.

A STABLE version is one which is both (a) working towards a milestone and (b) passes all
integration tests. It may be a bugfix, or a new feature. An UNSTABLE version is effectively
latest head, and is pushed on green - i.e. the unit tests and any internal integration tests
are passing. UNSTABLE versions will never make it to prod, only STABLE versions are pushed
to the prod channel.

Thus since we control the server API we expect that both (a) the API is largely stable once we
reach the first MILESTONE, and that (b) our clients will migrate off of deprecated API
elements over time and we can sunset them, rather than maintain older versions.

## Control

UNSTALBLE versions are effectively a pull request that is merged into main branch. Overseer then
runs the suite of integration tests over that version and increments the STABLE version if all
integration tests pass for that STABLE version. MILESTONES are a set of defined features that must
pass before overseer is able to increment the MILESTONE number.

Thus a MILESTONE is a set of features, a feature is embodied in integration tests than run in the cluster,
with each feature being necessarily stand alone. We support defining the feature before writing the
integration test, thus allowing the MILESTONE to develop organically.

A feature necessarily ties to a set of integration tests, where each must succeed in order for the STABLE
version to jump.
