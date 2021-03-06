# matts-tool
Matt's tool is a swiss-army tool of making life easier working with conjur.
It currently has a 2-in-1 tool for managing policy.

# Policy transformation
Given a set of v4 policy, it can convert that policy to valid v5 policy by stripping out unused
or unsupported items.  But it contains the actual reason for starting it...`!include` tags can be
included with the format `!include <filename>` and the file will be parsed and loaded and replace
the node it was included at.  This allows for use of v5 policy with `!include` tags for better
organization and handling of policy as well as being able to turn most v4 policy into a valid v5
policy document.

`matts-tool policy -i <filename>` will read in the given file, process any `!include` statements
(as well as processing `!include` statements in included files as well) and then remove any
incompatible/unnecessary tags, outputtig a valid v5 policy document that corresponds to what the
v4 policy represented.

`matts-tool policy -i <filename>` will also work for v5 policy that includes `!include` tags as well,
making organization and management of large policies reasonable.

# DB movement stuff
This will take a particular format pg_dump sql file from a v4 instance and load it into a v5 compatible
schema, decrypted and re-encrypting the values for password hashes and api keys.  Two values that are
necessary that are not command line options.  `IN_CONJUR_DATA_KEY` and `OUT_CONJUR_DATA_KEY`, with
`IN_CONJUR_DATA_KEY` being the key that was used to encrypt the values in the data dump file and
`OUT_CONJUR_DATA_KEY` the key to use for encrypting the values.

The following will generate a proper db dump that can be read and processed:

`pg_dump --data-only --schema="authn" --table="authn.users" > ~/data.sql`