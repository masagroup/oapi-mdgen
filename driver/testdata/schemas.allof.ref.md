
# Sample (1.0)
sample tests

## default


<a id="getTest"></a>
<details>
<summary><code>GET</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>


### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK | [Circle](#componentsSchemasCircle) |

</details>



## Schemas

<a id="componentsSchemasCircle"></a>
<details>
<summary><code>Circle</code></summary>

### All Of
> | schema |
> | --- |
> | [Shape](#componentsSchemasShape) |
> | object |

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | name | string | true | The shape name. |
> | color | string | true | The circle color, overrides the shape color. |
> | radius | number(double) | true | The circle radius. |


</details>

<a id="componentsSchemasShape"></a>
<details>
<summary><code>Shape</code></summary>

### Properties
> | name | type | required | description |
> | --- | --- | --- | ---  |
> | name | string | true | The shape name. |
> | color | string | false | The shape color. |


</details>



