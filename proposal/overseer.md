# Overseer

The goal of the overseer is to support a plugin type thing for handling
integration tests and validation of existing infrastructure and automatic
version promotion given validation passing.

## Process

Each project supports an optional validator which interacts with the overseer
in both (a) setting some basic configuration that defines some rules for
version promotion and (b) supports running validation tests through the
individual validators. Thus overseer is a go between and aggregator for the
validators.

## Key concepts

Validation supports versioned configuration - so we bascially say to overseer
that here is a list of validators, if you are able to run all of these then
then promote the github version and consider us done. We can then up the version
of the config, and we must add a new validation option in order to do that.

We can also add in a bug fix option that supports updating the minor version.

## Semantic versioning

Basically the versioning scheme is major . minor . patch

so effectively major changes are API breaking, minor changes support the
existing API, and patch changes are bug fixes placed in the queue. So each validation
set refers to the current API, and we update the API by adjusting the minor version once
the system is validated. To support a bug fix update we add a validator that captures
the particular bug, and then run against that. Otherwise our individual code pushes
are set against -alpha-1, -alpha-2 etc.

Thus the versioning is controlled by the overseer which will only promote alpha versions
once we add a new verisoning config. The versioning config figures out the promotion it
needs to make in advance and stores compliance / regression data accordingly.

## Configuration

The configuration passed to overseer defines which piece of the puzzle we'll update, we cannot
add a new config that updates major or minor until the existing one is done - but we can patch
an existing version if we find a bug fix.

So our config defines: the promotion type (MAJOR / MINOR / PATCH) and the set of validators that
must pass in order for overseer to make a promotion. We then validate against the currently running
version in the cluster. If we pass all validators, we promote, otherwise we mark and store
the failure. We can use the stored output to track the state of the system and the passing validators.

## Tasks

1.  
