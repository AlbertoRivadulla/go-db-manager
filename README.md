# Car maintenance tracker

Build the CLI app with
```bash
go build -o carmaintenance-cli cli/main.go
```
and run it using
```bash
./carmaintenance-cli -dbpath <database_path> -specs <specs_dir>
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
    table:
      name: MyTable
      schema: public

      foreign_keys:
        - column: other_id
          references:
            table: OtherTable
            column: id

      columns:
        - name: id
          type: integer
          auto_increment: true
          primary_key: true
          nullable: false

        - name: other
          type:integer
          nullable: false

        - name: title
          type: string

        - name: cost
          type: decimal(10,2)
          nullable: false
          default: "0.0"
          check: "price >= 0.0"
    ```
- `rules/`: defines rules for the regular maintenances. Example:
    ```YAML
    # TODO:
    ```
- `queries/`: specifies different queries to be run on the tables, with placeholders for values that can be supplied by the user. Example:
    ```YAML
    # TODO:
    ```


