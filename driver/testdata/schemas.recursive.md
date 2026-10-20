
# Sample (1.0)
sample tests

## default


<a id="getTest"></a>
<details>
<summary><code>GET</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>


### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK | [Node](#componentsSchemasNode) |

</details>



## Schemas

<a id="componentsSchemasNode"></a>
<details>
<summary><code>Node</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | name | string | true | The node name. |
> | parent | [Node](#componentsSchemasNode) | false |  |
> | children | [[Node](#componentsSchemasNode)] | false |  |


</details>



