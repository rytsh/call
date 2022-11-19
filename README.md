![call](./static/logo/call.svg)

Call's web page based on [Hugo](https://gohugo.io/).

## Development

First download the `hugo` cli.

```sh
brew install hugo
```

Install the front-end dependencies.

```sh
pnpm install
```

Now you can run the `hugo` server.

```sh
pnpm run dev
```

## Deployment

Deployment based on tag, add a tag with `v[0-9]*` prefix, then push it to the remote.  
Usually add `-doc` suffix to the tag name.

For build in local, run `pnpm build`.
