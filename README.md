# json2ts

json2ts converts JSON into valid Typescript classes with decorators.

## Installation

If you have go already installed you can simply run

```bash
go get github.com/sno6/json2ts
```

If you prefer a binary, you can also find those in the [releases section](https://www.github.com/sno6/json2ts/releases):

## Example

Take some JSON, such as this example from an HTTP response

```json
{
  "data": [{
    "type": "articles",
    "id": "1",
    "attributes": {
      "title": "JSON:API paints my bikeshed!",
      "body": "The shortest article. Ever.",
    },
    "relationships": {
      "author": {
        "data": {"id": "42", "type": "people"}
      }
    }
  }],
  "included": [
    {
      "type": "people",
      "id": "42",
      "included_attributes": {
        "name": "John",
        "age": 80,
        "gender": "male"
      }
    }
  ]
}

```

Copy the JSON and run the following:

```bash
>>> pbpaste | json2ts
```

Which will print the following Typescript code

```typescript
export class BaseClass {
	@IsDefined()
	@IsArray()
	@Type(() => Data)
	@ValidateNested({ each: true })
	public data!: Data[];

	@IsDefined()
	@IsArray()
	@Type(() => Included)
	@ValidateNested({ each: true })
	public included!: Included[];
}

export class Data {
	@IsDefined()
	@IsString()
	public type!: string;

	@IsDefined()
	@IsString()
	public id!: string;

	@IsDefined()
	@Type(() => Attributes)
	@ValidateNested({ each: true })
	public attributes!: Attributes;

	@IsDefined()
	@Type(() => Relationships)
	@ValidateNested({ each: true })
	public relationships!: Relationships;
}

export class Attributes {
	@IsDefined()
	@IsString()
	public title!: string;

	@IsDefined()
	@IsString()
	public body!: string;
}

export class Relationships {
	@IsDefined()
	@Type(() => Author)
	@ValidateNested({ each: true })
	public author!: Author;
}

export class Author {
	@IsDefined()
	@Type(() => Data)
	@ValidateNested({ each: true })
	public data!: Data;
}

export class Data {
	@IsDefined()
	@IsString()
	public id!: string;

	@IsDefined()
	@IsString()
	public type!: string;
}

export class Included {
	@IsDefined()
	@IsString()
	public type!: string;

	@IsDefined()
	@IsString()
	public id!: string;

	@IsDefined()
	@Type(() => IncludedAttributes)
	@ValidateNested({ each: true })
	public included_attributes!: IncludedAttributes;
}

export class IncludedAttributes {
	@IsDefined()
	@IsString()
	public name!: string;

	@IsDefined()
	@IsNumber()
	public age!: number;

	@IsDefined()
	@IsString()
	public gender!: string;
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