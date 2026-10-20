
# Sample (1.0)
sample tests

## default


<a id="postTest"></a>
<details>
<summary><code>POST</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>

### Parameters
> | name | in | data type | required | description |
> | --- | --- | --- | --- | --- |
> | body | body | [[Parameter](#componentsSchemasParameter)] | false | parameters  |

### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK |  |

</details>



## Schemas

<a id="componentsSchemasBooleanParameter"></a>
<details>
<summary><code>BooleanParameter</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | type | string<br/>*Available values:* boolean | true |  |
> | value | boolean | true |  |


</details>

<a id="componentsSchemasNullParameter"></a>
<details>
<summary><code>NullParameter</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | type | string<br/>*Available values:* null | true |  |


</details>

<a id="componentsSchemasParameter"></a>
<details>
<summary><code>Parameter</code></summary>

### One Of
Discriminated by `type`
> | value | schema |
> | --- | --- |
> | string | [StringParameter](#componentsSchemasStringParameter) |
> | boolean | [BooleanParameter](#componentsSchemasBooleanParameter) |
> | null | [NullParameter](#componentsSchemasNullParameter) |
> | UnmappedParameter | [UnmappedParameter](#componentsSchemasUnmappedParameter) |

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | type | string | true | Discriminator type for the parameter |


</details>

<a id="componentsSchemasStringParameter"></a>
<details>
<summary><code>StringParameter</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | type | string<br/>*Available values:* string | true |  |
> | value | string | true |  |


</details>

<a id="componentsSchemasUnmappedParameter"></a>
<details>
<summary><code>UnmappedParameter</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | type | string<br/>*Available values:* UnmappedParameter | true |  |


</details>



