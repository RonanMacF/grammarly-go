# nvim-grammarly

This aims to provide means of communicating to grammarly via commandline to lint english in files.

## Configuration

A configuration file can be passed to provide specific error reporting for your needs. These configuration parameters wlll sent to Grammarly when initiaitng the socket communication. By default a config file will be loaded from `~/.config/grammarly-go/config.toml`. A Custo path for a config file can be provided by bayying the `--configFile` flag followed by a path. If a config file is not found sensible defaults will be used which can be seen in the sample configuration file in the repo.

## Integration with editors

A strong usecase for me is using this tool to lint errors in my editor. This can be easily accomplished by using the editor agnostic tool efm-langserver. For example with Neovim you can achieve linting in text files using the following configuration with efm
```
require"lspconfig".efm.setup {
    filetypes = {"text"}
    settings = {
        rootMarkers = {".git/"},
        languages = {
            text = {
               lintCommand = "grammarly-go -filePath ",
               lintIgnoreExitCode = true,
               lintFormats = {"%f:%l:%c: %m"},
            }
        }
    }
}
```

## Contribution
> make it work, 
make it right, 
make it fast

This project is deep in the first stage now. PR's, feature requests, bug reports are all welcome. I have an idea of what I want, if you know what you want let me know and I can try to implement it.

## Note
There is another language server out there which is probably more feature rich than this. There is a dedicated [VSCode extension](https://github.com/znck/grammarly), there is a dedicated [emacs plugin](https://github.com/emacs-grammarly/unofficial-grammarly-language-server) and a dedicaed [COC plugin](https://github.com/gianarb/coc-grammarly).  
