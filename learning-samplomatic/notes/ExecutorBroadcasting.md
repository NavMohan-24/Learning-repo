
The input and output data to Executors are array. Broadcasting rules could be applied to control the input and output. 

- Intrinsic shapes --> The shape the data type would have by data. Depends upon the data.

  - For instance a parameterised circuit with 3 parameters have an intrinsic axis of 3 --> represented as (3,).

- Extrinsic shape --> Sweep dimensions. How many configuration you want to run. 

  - If the same circuit is executed against 5 parameter sets, then we represent the TOTAL SHAPE as (5, 3). Here, 5 is the extrinsic and 3 is the instrinsic axis. 
  - Extrinsic axis is the left most one & Intrinsice axis and right most one 

- While performing different process, the input arrays are brodcasted along extrinsic axis.

    - For instance, consider the same circuit which is now have to executed against 4 noise values.
    - The process of parameterisation results in a input array of (5,3) which we can rewrite as (5,1,3).
    - Now since, we need to execute for 4 different noise values which can be represented as (4,).
    - Now, combing both process the input array are broadcasted as (5,4,3)

- The shape of output array is determined as follows:
  
    - Extrinsic axis --> brodcasted input dimension
    - Intrinsic axis --> Determined by output type. For a sampler type of output, the instrinsic axis would be (num_shots,creg_size)

- Thus, for the above example of running a parameterised circuit (with 3 parameters) against 5 parameter set on 4 different noise levels the output array depends upon the number of classical registers and shots. Let's say if creg_size is 2 and num_shots = 1024. Then the shape of the output array (5,4,1024,2)

- While performing twirling or any other randomized sampling based experiment, the number of randomizations can be performed specified as a shape parameter.

  - For instance, consider the parameterized circuit with number of parameters = 3 executed against 5 parameter set. Thus, implicit input shape is (5,3). Now, for each parameter set , if we want to perform twirling with num_randomizations 10 for each parameter set we can specific it as the explict axis of shape parameter of `SamplexItem`.


```python 
from qiskit_ibm_runtime.quantum_program import QuantumProgram

program.append_samplex_item(circuit=template_circuit, 
                            samplex_arguments= {
                                "parameter_values": np.random.rand(5, 3)
                                }, shape = (10,5),) # 
    
```

We can also specific the randomization via multiple randomization axis.

```python 
from qiskit_ibm_runtime.quantum_program import QuantumProgram

program.append_samplex_item(circuit=template_circuit, 
                            samplex_arguments= {
                                "parameter_values": np.random.rand(5, 3)
                                }, shape = (2,5,5),) # 
    
```
This is useful is structured analysis?
