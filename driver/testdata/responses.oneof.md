
# Sample (1.0)
sample tests

## default


<a id="getTest"></a>
<details>
<summary><code>GET</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>


### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK | one of [Unit](#componentsSchemasUnit), [Party](#componentsSchemasParty) |
> | 400 | Bad Request | [one of [Unit](#componentsSchemasUnit), [Party](#componentsSchemasParty)] |

</details>



## Schemas

<a id="componentsSchemasParty"></a>
<details>
<summary><code>Party</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | kind | string | false |  |
> | name | string | false |  |


</details>

<a id="componentsSchemasUnit"></a>
<details>
<summary><code>Unit</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | kind | string | false |  |
> | id | integer | false |  |


</details>



