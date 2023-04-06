# kubepost

<img align="right" alt="kubepost" width="180px" src="assets/gopher.png">

<p>
    <a href="https://github.com/orbatschow/kontext/actions/workflows/default.yaml" target="_blank" rel="noopener"><img src="https://img.shields.io/github/actions/workflow/status/orbatschow/kontext/default.yaml" alt="build" /></a>
    <a href="https://github.com/orbatschow/kontext/releases" target="_blank" rel="noopener"><img src="https://img.shields.io/github/release/orbatschow/kontext.svg" alt="Latest releases" /></a>
    <a href="https://github.com/orbatschow/kontext/blob/master/LICENSE" target="_blank" rel="noopener"><img src="https://img.shields.io/github/license/orbatschow/kontext" /></a>
</p>

Kontext is a single binary application, that takes yet another approach on kubeconfig management.

## Features

Kontext has several features, that will ease your life when dealing with different kubeconfig files:

### Context

Switch between a context by just calling the binary, without any arguments. It will read your current kubeconfig file
and list all available options. To get more information about setting and getting a context run:

```shell
kontext get context --help
kontext set context --help
```

### Groups

Groups refer to one or more sources and can be used to bundle kubeconfig files together. You
can switch between groups and enable or disable multiple sources at once. Get help via:

```shell
kontext get group --help
kontext set group --help
```

### Sources

Source include or exclude kubeconfig files as a glob pattern. A source always computes all
included sources files first and then removes all duplicates. After the include section has
been computed the same happens for all files, that shall be excluded. Take a look at the 
[example](./example/kontext.yaml) to understand sources in depth.

## Installation

At the moment kontext is distributed as a single binary and can be downloaded from the 
[releases](https://github.com/orbatschow/kontext/releases).

## Configuration

Have a look at the [example](./example/kontext.yaml) file. It should be well described and show you
how to configure kontext.


## Contributing

Contributions are always welcome, have a look at the [contributing](docs/contributing.md) guidelines to get started.