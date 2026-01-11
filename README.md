# Car maintenance tracker

Build the CLI app with
```bash
go build -o cli-carmaintenance cli/main.go
```
and run it using
```bash
carmaintenance-cli -dbpath <database_path> -specs <specs_dir>
```
where `<specs_dir>` is a directory with YAML files that specify the different tables in the database, and queries to run on them. This contains the following subdirectories and files (there can be more than one file in each subdirectory):
```
specs/
├-- queries/
|   └-- queries.yaml
├-- rules/
|   └-- rules.yaml
└-- tables/
    └-- table.yaml
```
The files in these directories have the following contents:
- `tables/`: defines the different tables in the database, and requirements for the data in them. An example of a file located there would be:
    ```YAML
    # TODO:
    ```
- `rules/`: defines rules for the regular maintenances. Example:
    ```YAML
    # TODO:
    ```
- `queries/`: specifies different queries to be run on the tables, with placeholders for values that can be supplied by the user. Example:
    ```YAML
    # TODO:
    ```


