# GVM: Go Version Manager for Windows

<!--- These are examples. See https://shields.io for others or to customize this set of shields. You might want to include dependencies, project status and licence info here --->

![GitHub repo size](https://img.shields.io/github/repo-size/nicewook/gvm)
![GitHub contributors](https://img.shields.io/github/contributors/nicewook/gvm)
![GitHub stars](https://img.shields.io/github/stars/nicewook/gvm?style=social)
![GitHub forks](https://img.shields.io/github/forks/nicewook/gvm?style=social)

Project name is a `gvm` that allows `Go developers who are using Windows` to use multiple go version easily.
There is already `gvm`(https://github.com/moovweb/gvm) for Linux and MacOS, but I couldn't find
well working `gvm` for Windows OS, so the project started.

## Installing `gvm`

To install `gvm` on Windows, follow these steps:

```
$ go get -u github.com/nicewook/gvm
```

## Features

- version
- list
- listall
- install
- uninstall
- use
  

## Using `gvm` commands

### version

`version` command shows the gvm version.
```
$ gvm version 
```

### list

`list` command shows the installed go SDK(s), and `system` means originally installed go.


```
$ gvm list 
```

### listall

`listall` command shows the installed go SDK(s), and `system` means originally installed go.
```
$ gvm listall 
```

### install

`install` command installs go SDK(s)
- It checks if the version(s) already installed, or not existing version to install.
- It can install multiple versions at once.

```
$ gvm install 1.13.1 1.13.2
```

### uninstall

`uninstall` command uninstalls go SDK(s).
- It checks if the version(s) is/are not existing.
- It can uninstall multiple versions at once.

```
$ gvm uninstall 1.13.1 1.13.2
```

### use

`use` command changes using go SDK versions to desired.
- If no version specified, it shows the current using version of the go SDK
- If checks if the version is not installed.
- If you name the version as `system`, it will changes to the originally installed go version.


```
$ gvm use

$ gvm use system

$ gvm use 1.13.2
```

## Contact

If you want to contact me you can reach me at <nicewook@hotmail.com>

## License

<!--- If you're not sure which open license to use see https://choosealicense.com/--->

This project uses the following license: [<license_name>](link).

## Reference

- README.md template: https://github.com/scottydocs/README-template.md
- moovweb/gvm: https://github.com/moovweb/gvm
- Download go verions(Bill Kennedy): https://www.ardanlabs.com/blog/2020/04/modules-06-vendoring.html
```
