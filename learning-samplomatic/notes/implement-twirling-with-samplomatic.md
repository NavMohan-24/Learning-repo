## How Samplomatic performs twirling?

Consider a bell circuit to which we want to apply twirl operation on the two-qubit operation and measurement.  In samplomatic, we first need to group the operations using box and annotate them with `Twirl` directive. Samplomatic implements random Pauli gates for twirling using dressing. Thus, position where the dressing needs to be applied should be mentioned in the `Twirl` annotations. At this point, 
we have a annotated boxed up circuit (for simplicity we just call it as "boxed up circuit") which is considered as declarative model of randomization.  The boxed up circuit does not actually implements randomization. 

<div style = "height: 200px; width:auto">
<img src ="./images/demo-circuit-1.png" alt="boxed-up-circuit">
</div>

Samplomatic would take such a boxed up circuit as input and builds a procedural representation of randomization using `build` method. To have a procedural representation, the `build` method would output two objects; a Template circuit and a Samplex.

_Template circuit from the Bell circuit:_
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


The virtual gates generated at `R0` propagates leftwards and rightwards. Those propagated rightwards will get combined with the virtual gates propagates leftwards from `R1` and eventually gets collected in the dressing. Those propagates leftwards `R1` in their journey to get collected in the dressing, need to first move past through CX gate and then need to take care of the hadarmard operation in the first qubit (refer boxed up circuit). At the end of their jour, the virtual gates are converted to appropriate parameters such that the gates could be implemented using dressing. 

When a Pauli gate $P$ propagate across a clifford gate $C$, the Pauli gate would get transformed as $P^{'} = CPC^{\dagger}$. Thus, propagating across `CX` gate modifies or mutates the virtual gates. Similarly, Hadamard gate is taken care by right multiplying it against the virtual gates on the qubit-0. In short, virtual gates propagated leftwards will get mutated before getting collected in the dressing. Hence, the parameters in the dressing should be modified appropriately during execution. Now, the gates propagated rightwards from `R1` will be applied as bit flips on final measurement results during postprocessing.

In short, dressing is used to implement both virtual gates and the gates present in the original. However, the template circuit alone are unaware of the the parameters to implement. Further, we also require the info about the bit flips that should be applied to measurement results of template the circuit to undo the effect of measurement twirling. Samplomatic uses `samplex` to supply parameters and the array values corresponding to bit flips required during execution of template circuit.

### Samplex 

Samplex is a core-type defined in samplomatic. It represents a probability distribution of the parameter values for executing template circuits and classical quantities to perform the post-processing. Samplex encodes the process of randomization as graph based procedural representation. Running `samplex.graph` would return a Directed Acyclic Graph (DAG) in which each node represents a process.

_Samplex from the Bell circuit:_
<div style = "height: 300px; width:600px;">
<img src ="./images/samplex.png" alt="samplex" style="width: 100%; height: 100%; object-fit: cover;">
</div>

In a samplex, there are three types of nodes:

- Sampling Nodes (Star Shaped) 
- Intermediate Nodes (Circles)
- Collection Nodes (Bow ties)

The sampling nodes are responsible for instantiating virtual registers. Depending upon task in hand, whether it is generating randomizations for twirling, sampling noise from noise models or injecting basis, these nodes performs the sampling of virtual gates from a set. The intermediate nodes represents all the mutations that are happening to the  virtual register as they propagate across the circuit. The mutations can be of different types, such as:

   - Combining virtual register to one,
   - Commute virtual gates across other gate operations,
   - Change representation for eg; convert a pauli operator to equivalent U2 gate representation.

The collection nodes are resposible to convert virtual gates to parameter values (blue bow tie) for template circuits or other array valued fields for post-processing.

As we discussed, the samplex represents a probability distribution and the interface (API) to draw samples from it are `samplex.sample()`. It will return a collection of arrays that will be used during execution of template circuits .However, for the samplex to return outputs we need to bind all the input values that required. For the task of twirling bell circuit there are no inputs of required. However, inputs can arise if the original circuit is parameteric. Further, if the dressing are Pauli Linblad Maps for noise injection or the basis change array, inputs are required. The inputs and outputs of samplex are strongly typed, meaning for any particular instance of a samplex, their names, types, and shapes are fixed and queryable before any sampling is performed. The easiest way to see the required inputs and expected outputs is to print the Samplex object.




<!-- Appropriate parameters should be supplied by taking into consideration of the mutations happening to the circuit and gate operations present in the original circuit. However, the template circuit cannot keep track of the mutations happening to virtual gates. In order to keep track of this, samplomatic uses a special built in type called Samplex. The samplex is a graph based procedural representation of randomization process. The samplex could be depicted as a Direct Acyclic Graph, where each node represent a procedure.

It represent a probability distribution 


The samplex could be represented as a DAG which represents the process of randomization.  

 The virtual gates propagating leftwards will get mutated as it propagates through $CZ$ gate due to commutation relations. Now, the Hadamard operation on the qubit-0 also have to be taken care. The virtual gates at qubit-0 will be right multiplied against Hadmard gates. Then virtual gates needs to be converted into parameter values and will get collected in the dressing (between `L0` and `M0`). Similarly, virtual gates that propagates rightwards from `R0` also get accumulated in the dressing between `L1` and `M1`.

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
