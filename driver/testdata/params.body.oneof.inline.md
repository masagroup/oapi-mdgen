
# Sample (1.0)
sample tests

## default


<a id="postTest"></a>
<details>
<summary><code>POST</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>

### Parameters
> | name | in | data type | required | description |
> | --- | --- | --- | --- | --- |
> | body | body | object | false | Creation request body  |
> | └── target | body | one of [Unit](#componentsSchemasUnit), [Party](#componentsSchemasParty) | true |  |
> | └── value | body | any of string, integer | false |  |

### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK |  |

</details>



## Schemas

<a id="componentsSchemasParty"></a>
<details>
<summary><code>Party</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | id | integer | false |  |


</details>

<a id="componentsSchemasUnit"></a>
<details>
<summary><code>Unit</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | id | integer | false |  |


</details>



