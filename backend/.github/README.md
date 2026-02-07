# Hexagonal

#### Summary

1. [Description](#description-)
2. [Dependencies](#dependencies-)
3. [Make Command](#make-file-)
4. [Running the project](#running-the-project-)
   1. [Repository cloned without Go installation](#repository-cloned-without-go-installation)
   2. [Repository cloned with Go installation](#repository-cloned-with-go-installation)
   3. [From built version](#from-built-version)
5. [Stopping the project](#stopping-the-project-)
   1. [Maintaining the data](#maintaining-the-data)
   2. [Clearing the data](#clearing-the-data)

<h1></h1>

1. #### Description [&uarr;](#summary)

   Hexagonal architecture written in Go.

2. #### Dependencies [&uarr;](#summary)
   - Make
   - Docker & Docker Compose
   - Go 1.25+ (To Dev.)

3. #### Make File [&uarr;](#summary)

   <details>
   <summary>Commands:</summary>

   ```sh
    Usage:
      make [COMMAND]

    Example:
      make help

    Commands:

    help                           Display available commands and their descriptions
    init                           Create environment file
    test                           Run tests and generate coverage report
    run                            Run application from source code
    build                          Build the all applications from source code
    swag                           Update swagger files
    compose-up                     Create and start containers
    compose-build                  Build, create and start containers
    compose-down                   Stop and remove containers and networks
    compose-remove                 Stop and remove containers, networks and volumes
    compose-exec                   Access container bash
    compose-log                    Show container logger
    compose-top                    Display containers processes
    compose-stats                  Display containers stats
    go-benchmark                   Benchmark code performance
    go-lint                        Run lint checks
    go-audit                       Conduct quality checks
    go-format                      Fix code format issues
    go-tidy                        Clean and tidy dependencies
   ```

   </details>

4. #### Running the project [&uarr;](#summary)
   1. ##### Repository cloned without Go installation

      Follow these steps if you cloned the repository and do not have golang instalation:
      1. Open the terminal in the cloned repository folder.
      2. Run `make init compose-build base=source` to create environment file and create and start production containers.

   2. ##### Repository cloned with Go installation

      Follow these steps if you cloned the repository and have golang instalation:
      1. Open the terminal in the cloned repository folder.
      2. Run `make init build compose-build` to create environment file and create and start production containers.

   3. ##### From built version

      Follow these steps if you downloaded the built version:
      1. Open terminal in built version folder.
      2. Run `make compose-build` to create and start production containers.

5. #### Stopping the project [&uarr;](#summary)
   1. ##### Maintaining the data

      Follow these steps if you want to stop the project but don't want to delete all the data:
      1. Open the terminal in the project folder.
      2. Run `make compose-down`

   2. ##### Clearing the data

      Follow these steps if you wish to stop the project and delete all the data:
      1. Open the terminal in the project folder.
      2. Run `make compose-remove`

<p style="text-align:right">&#40;<a href="#title">back to top</a>&#41;</p>
