
# Sample (1.0)
sample tests

## default


<a id="getTest"></a>
<details>
<summary><code>GET</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>


### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK | [Pet](#componentsSchemasPet) |

</details>



## Schemas

<a id="componentsSchemasCat"></a>
<details>
<summary><code>Cat</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | meow | boolean | true | Does it meow. |


</details>

<a id="componentsSchemasDog"></a>
<details>
<summary><code>Dog</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | bark | boolean | false | Does it bark. |


</details>

<a id="componentsSchemasPet"></a>
<details>
<summary><code>Pet</code></summary>

### One Of
> | schema |
> | --- |
> | [Cat](#componentsSchemasCat) |
> | [Dog](#componentsSchemasDog) |
> | object |

</details>



