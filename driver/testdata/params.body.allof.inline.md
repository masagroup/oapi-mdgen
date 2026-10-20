
# Sample (1.0)
sample tests

## default


<a id="postShapes"></a>
<details>
<summary><code>POST</code> <code><b>/shapes</b></code> <code><i>creates a shape</i></code></summary>

### Parameters
> | name | in | data type | required | description |
> | --- | --- | --- | --- | --- |
> | body | body | object | false | Creation request body  |
> | └── name | body | string | true |  |
> | └── uuid | body | string | false |  |
> | └── category | body | string | true |  |
> | └── points | body | [object] | false |  |
> | └── origin | body | object | false |  |
> | &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└── x | body | integer | false |  |
> | &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└── y | body | integer | false |  |

### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK |  |

</details>

<a id="putShapes"></a>
<details>
<summary><code>PUT</code> <code><b>/shapes</b></code> <code><i>updates a shape</i></code></summary>

### Parameters
> | name | in | data type | required | description |
> | --- | --- | --- | --- | --- |
> | body | body | <dl><dt>**application/json**</dt><dd>[object](#putShapesBodyApplicationJson)</dd><dt>**application/xml**</dt><dd>[Shape](#componentsSchemasShape)</dd></dl> | false | Update request body  |
### Body Schema

<a id="putShapesBodyApplicationJson"></a>
object
> | name | type | required |description |
> | --- | --- | --- | --- |
> | └── name | string | true |  |
> | └── uuid | string | false |  |
> | └── category | string | false |  |
### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK |  |

</details>



## Schemas

<a id="componentsSchemasShape"></a>
<details>
<summary><code>Shape</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | name | string | true |  |
> | uuid | string | false |  |


</details>



