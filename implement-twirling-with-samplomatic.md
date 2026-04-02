## How Samplomatic performs twirling?

Consider a bell circuit to which we want to apply twirl operation on the two-qubit operation and measurement.  In samplomatic, we first need to group the operations using box and annotate them with `Twirl` directive. Samplomatic implements random Pauli gates for twirling using dressing. Thus, position where the dressing needs to be applied should be mentioned in the `Twirl` annotations. At this point, 
we have a annotated boxed up circuit (for simplicity we just call it as "boxed up circuit") which is considered as declarative model of randomization.  The boxed up circuit does not actually implements randomization. 


<div style = "height: 200px; width:auto">
<img src ="./images/demo-circuit-1.png" alt="boxed-up-circuit">
</div>

<!-- ![Boxed up Bell circuit](./images/demo-circuit-1.png){height=20%} -->



Samplomatic would take such a boxed up circuit as input and builds a procedural representation of randomization using `build` method. To have a procedural representation, the `build` method would output two objects; a Template circuit and a Samplex.

Template circuit from the Bell circuit:

![Template circuit](./images/template-circuit-1.png)



<!-- ![Samplex of Bell circuit](./images/samplex.png){height=40%} -->

The template circuit will be structurally similar to boxed up circuit. They use barriers instead of Boxes for the logical partioning. In the template circuit, the barriers `L0` and `R0` represents the scope of Box-1 and `L1` and `R1` represents the scope of Box-2. The parametric gates between `L0` (`L1`) and `M0`(`M1`) is the dressing and by default it will be on the left side of the box unless specified. The template circuit will be executed inplace of the boxed up circuit as required randomizations could easily implemented choosing appropriate parameter values for dressing.


<!-- For instance, the span between `L0` (`L1`) and `R0` (`R1`) denotes the first (second) box. The parametric gates between `L0` (`L1`) and `M0`(`M1`) are the dressing and by default it will be on the left side of the box unless specified. -->
<!-- Further, parameter values The parameter values required for randomization will be supplied by Samplex object, which is a main type defined in Samplomatic. <!-- Using the boxed up circuit with annotations, samplomatic construct an template circuit which will be used for randomizations in the subsequent steps.  The template circuit of the above circuit is as follows: -->
<!-- The dressing gates are used to implement gate operations in the circuit such as `Hadamard` operation in the qubit-0 and random Pauli gate operations that are applied due to twirling. As the name implies, the template circuit defines the structure of the circuit that will be executed. 
The template circuit does not store or remembers the Hadamard operation on qubit-0 within Box-1. It is the responsibility of Samplex object to implement the Hadamard gate along with the Pauli gates for twirling during runtime. Before, delving into Samplex, let's look into how twirling will get implemented in Samplomatic. -->

#### Implementating twirling with dressing

The random Pauli gates used for twirling are considered `virtual` because they do not add any additional operations to the circuit. Instead, they act as a directive to alter how adjacent single-qubit gates are implemented. Samplomatic generates virtual gates (random Pauli gates) on the boundary of the box opposite to its dressing. Specifically, the virtual gates are generated at `R0` and `R1` of the template circuit and propagated leftward and rightward, ultimately getting accumulated in the dressing.

<!-- In the case of template circuit we are discussing, virtual gates are generated at the barriers are generated at `R0` and `R1` and propagated left and rightwards.  -->

![marked-dressed-circuit](./images/marked-template-circuit-1.png)

The journey of virtual gates generated at `R0` and `R1` is depicted below.

![samplex-depiction](./images/pre-samplex.svg)

( change blue boxes --> convert virtual gates to parameter values.)

The virtual gates generated at `R0` propagates leftwards and rightwards. Those propagated rightwards will get combined with the virtual gates propagates leftwards from `R1` and eventually gets collected in the dressing. Those propagates leftwards `R1` in their journey to get collected in the dressing, need to first move past through CX gate and then need to take care of the hadarmard operation in the first qubit (refer boxed up circuit).  

When a Pauli gate $P$ propagate across a clifford gate $C$, the Pauli gate would get transformed as $P^{'} = CPC^{\dagger}$. Thus, propagating across `CX` gate modifies or mutates the virtual gates. Similarly, Hadamard gate is taken care by right multiplying it against the virtual gates on the qubit-0. In short, virtual gates propagated leftwards will get mutated before getting collected in the dressing. Hence, the parameters in the dressing should be modified appropriately during execution. Now, the gates propagated rightwards from `R1` will be applied as bit flips on final measurement results during postprocessing.

In short, dressing is used to implement both virtual gates and the gates present in the original circuit. Appropriate parameters should be supplied by taking into consideration of the mutations happening to the circuit and gate operations present in the original circuit. However, the template circuit cannot keep track of the mutations happening to virtual gates. In order to keep track of this, samplomatic uses a special built in type called Samplex. The samplex is a graph based procedural representation of randomization process. The samplex could be depicted as a Direct Acyclic Graph, where each node represent a procedure.

Samplex from the Bell circuit:
<div style = "height: 300px; width:600px;">
<img src ="./images/samplex.png" alt="samplex" style="width: 100%; height: 100%; object-fit: cover;">
</div>

It represent a probability distribution 


The samplex could be represented as a DAG which represents the process of randomization. 


<!-- The virtual gates propagating leftwards will get mutated as it propagates through $CZ$ gate due to commutation relations. Now, the Hadamard operation on the qubit-0 also have to be taken care. The virtual gates at qubit-0 will be right multiplied against Hadmard gates. Then virtual gates needs to be converted into parameter values and will get collected in the dressing (between `L0` and `M0`). Similarly, virtual gates that propagates rightwards from `R0` also get accumulated in the dressing between `L1` and `M1`.

Similar procedure is also followed at the barrier `R1` and virtual gate layers $Q \cdot Q$ will be generated. The virtual gates propagating leftwards will be collected in the dressing between `L1` and `M1`.  -->
<!-- In short,  samplomatic converts a boxed up circuits with directives into a procedural representation of the randomization. The template circuit gives the structure of the circuit that have to executed. The samplex will keep track of the mutation happening to virtual gates as it propagates across the circuit and gives the parameter values for dressing layer.  -->






> [!NOTE]
> The Propagation of Pauli gates across of Clifford gates follows this logic: 
>
> $$C \cdot P = C \cdot P \cdot (C^{\dagger}C) 
= (CPC^{\dagger})\cdot C = P^{'}\cdot C$$ 
> The above expression explains the logic of propagating **from right to left**. While propagating gates **from left to right**:
> $$ P \cdot C = (CC^{\dagger}) \cdot P \cdot C = C \cdot (C^{\dagger}PC) = C \cdot P^{'} $$







> [!Caution]
> . How Pauli gates propagates across non-clifford gates?
>
> . Can dressing have two qubit operations like Rzz($\theta$)? What happens if we box up an Rzz operation with CZ operation and apply twirling?
> 
> . Why gates are not getting mutated as it move past the measurement operations?
> 
> . Is it possible to do
