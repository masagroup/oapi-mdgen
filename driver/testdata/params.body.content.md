
# Sample (1.0)
sample tests

## default


<a id="postTest"></a>
<details>
<summary><code>POST</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>

this is a sample test

### Parameters
> | name | in | data type | required | description |
> | --- | --- | --- | --- | --- |
> | body | body | <dl><dt>**application/xml**</dt><dd>[ObjectType1](#postTestBodyApplicationXmlObjectType1)</dd><dt>**application/json**</dt><dd>[object](#postTestBodyApplicationJson)</dd></dl> | false |   |
### Body Schema

<a id="postTestBodyApplicationXmlObjectType1"></a>
ObjectType1
> | name | type | required |description |
> | --- | --- | --- | --- |
> | └── param1 | string | false |  |
> | └── param2 | string | false |  |
<a id="postTestBodyApplicationJson"></a>
object
> | name | type | required |description |
> | --- | --- | --- | --- |
> | └── param1 | string | false |  |
> | └── param2 | string | false |  |
</details>





