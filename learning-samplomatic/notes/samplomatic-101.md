# Samplomatic - 101



Box operations and annotations are the corner stone of the Samplomatic. 

## Box Operation
- Control Flow construct (like `if_test`, `switch` etc.) that can be added to circuits with an explicit condition.
- Intuitively, it acts like a logical partition such that all operations inside the box operation are treated as atomic.

   - Contents inside the box operations behave like as if the there are barriers at the start and end of the box operation.
   - However, unlike normal barriers, it is permissable to commute operations with box operations if all the gates/operations in the box commutes with the operation outside the box.

- Box can also be used to define an explicit scope to circuit elements like `Variables`, `Stretches` and also compiler/transpiler passes.

### Implementation

A box operation can be implemented in two ways:

- _The `BoxOp` way_ : Passing `QuantumCircuit` along with `qubits` and `clbits` to `BoxOp` and appending it to the target circuit.
```python
from qiskit.circuit import QuantumCircuit, BoxOp
 
body_0 = QuantumCircuit(4)
body_0.cz(0, 1)
body_0.cz(2, 3)

qc = QuantumCircuit(9)
qc.box(body_0, [0, 1, 2, 3], [])
```

- _Builder-Interface Form_ : use other `QuantumCircuit` methods within the Python `with` scope to add instructions to the `box`.

```python
from qiskit.circuit import QuantumCircuit

qc = QuantumCircuit(4)
with qc.box():
    qc.h(0)
    qc.cz(0, 1)
    qc.cz(2, 3)

qc.measure_all()
```


Both approach would give the following circuit:
```
     ┌───────     ───────┐ 
q_0: ┤        ─■─        ├─
     │         │         │ 
q_1: ┤        ─■─        ├─
     │ Box-0       End-0 │ 
q_2: ┤        ─■─        ├─
     │         │         │ 
q_3: ┤        ─■─        ├─
     └───────     ───────┘ 
q_4: ──────────────────────
                           
q_5: ──────────────────────
                           
q_6: ──────────────────────
                           
q_7: ──────────────────────
                           
q_8: ──────────────────────
```

### Annotations

Annotations is a framework used to attach metadata to Box operations within quantum circuits and DAG circuits. The metadata added will be consumed by any transpiler passes during transpilation.

They could be compared to the `PropertySet` object which is a shared dictionary that is used by transpiler passes during transpilation. A typical usage pattern is as follows, an analysis pass analyzes the properties of the `DAGCircuit` and stores it in `PropertySet` and a transformation pass reads the `PropertySet` and use it to modify the `DAGCircuit`.

The scope of `Annotations` is very local compared to `PropertySet`. Usually it is applyed to a box of instructions. It might be also present in the output of transpile function, if the it is intended for further consumption by a lower-level part of your backend’s execution machinery. 

> For example, an annotation might include metadata instructing an error-mitigation routine to treat a particular box in a special way

Annotations could be added to box as follows:

```python
from qiskit.circuit import QuantumCircuit, Annotation

class MyAnnotation(Annotation):
    namespace = "my.namespace"

qc = QuantumCircuit(9)
with qc.box([MyAnnotation()]):
    qc.cz(0, 1)
    qc.cz(2, 3)
```

> Note: ATM, the passmanager does not check the validity of Annotations. Thus it is possible that a custom Annotation could get invalidated by certain transformation passes. It is user's response to ensure that the that the compiler passes selected will not invalidate the annotation.


Docs: https://quantum.cloud.ibm.com/docs/en/api/qiskit/qiskit.circuit.QuantumCircuit#box


## Boxes and Annotations in Samplomatic

Boxes and Annoations are new features in Qiskit, which allow Samplomatic to perform randomization of quantum circuits which is crucial in noise learning and tailoring tasks. Boxes defines the scope of an operation applied the qubits amd annotations defines the _directives_ and _dressing_ that we apply to box operations. 

- Directives --> Defines the action to be done with the boxes.
- Dressing --> Defines the group of parameterized gate operations that are added to left or right side of the box. 


Some of the important directives supported by Samplomatic:

1. Twirl - Directive to twirl the contents of box instruction.
2. ChangeBasis - Directive to add basis changing gates.
3. InjectNoise - Directive to inject noise into the box instruction.

The dressing could consist of gate operations from the box, as well as operations required to enact directives. 

- Gates in the box that are compatible with and on the same side of parametrized gates
- Tirling gates, noise injection and basis changing gates.


Dressing and simultaneous two qubit entangling operation together composes the **Dressed Layer** which is the basic building block any noise-learning or error mitigation/suppression circuits (application circuit).

![dressed-layer](./images/dressed_layer.png)




<!-- The simplest case of dressing layer is the random Pauli gates that applied both sides of a two-qubit gate during twirling. In Samplomatic, twirling of a two-qubit gate could be done as follows: -->

<!-- ```python
from qiskit.circuit import QuantumCircuit
from samplomatic import Twirl

qc = QuantumCircuit(2)
with qc.box(annotations=Twirl())

``` -->