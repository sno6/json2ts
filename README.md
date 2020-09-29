# json2ts

json2ts converts JSON into valid Typescript classes with decorators.

## Installation

If you have go already installed you can simply run

```bash
go get github.com/sno6/json2ts
```

If you prefer a binary, you can also find those in the [releases section](https://github.com/sno6/json2ts/releases/tag/v1.0)

## Example

Take some JSON, such as this example from an HTTP response

```json
{
  "data": [{
	  "status": "success",
	  "flag": true,
	  "array": [{
		  "ofData": "Also supports any level of depth!"
	  }]
  }]
}
```

Pass the JSON through json2ts

```bash
# Either from the clipboard (OSX)
>>> pbpaste | json2ts -d

# Or from a file
>>> json2ts -i my_file.json -d
```

and voila, your Typescript generated code

```typescript
export class BaseClass {
	@IsDefined()
	@IsArray()
	@Type(() => Data)
	@ValidateNested({ each: true })
	public data!: Data[];
}

export class Data {
	@IsDefined()
	@IsString()
	public status!: string;

	@IsDefined()
	@IsBoolean()
	public flag!: boolean;

	@IsDefined()
	@IsArray()
	@Type(() => Array)
	@ValidateNested({ each: true })
	public array!: Array[];
}

export class Array {
	@IsDefined()
	@IsString()
	public ofData!: string;
}
```

## Usage

```bash
Transform JSON into typescript classes

Usage:
  json2ts [flags]

Flags:
  -d, --decorators      Add decorators to class parameters
  -i, --input string    Optional input file
  -o, --output string   Optional output file (defaults to stdout)
  -p, --prefix string   Prefix for all class names
  -r, --root string     Name of the root class
  -h, --help            help for json2ts
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
