??
?#?#
W
AddN
inputs"T*N
sum"T"
Nint(0"!
Ttype:
2	??
D
AddV2
x"T
y"T
z"T"
Ttype:
2	??
E
AssignAddVariableOp
resource
value"dtype"
dtypetype?
B
AssignVariableOp
resource
value"dtype"
dtypetype?
~
BiasAdd

value"T	
bias"T
output"T" 
Ttype:
2	"-
data_formatstringNHWC:
NHWCNCHW
~
BiasAddGrad
out_backprop"T
output"T" 
Ttype:
2	"-
data_formatstringNHWC:
NHWCNCHW
R
BroadcastGradientArgs
s0"T
s1"T
r0"T
r1"T"
Ttype0:
2	
N
Cast	
x"SrcT	
y"DstT"
SrcTtype"
DstTtype"
Truncatebool( 
8
Const
output"dtype"
valuetensor"
dtypetype
9
DivNoNan
x"T
y"T
z"T"
Ttype:

2
A
FloorDiv
x"T
y"T
z"T"
Ttype:
2	
B
GreaterEqual
x"T
y"T
z
"
Ttype:
2	
.
Identity

input"T
output"T"	
Ttype
9
	IdentityN

input2T
output2T"
T
list(type)(0
?
	LessEqual
x"T
y"T
z
"
Ttype:
2	
,
Log
x"T
y"T"
Ttype:

2
q
MatMul
a"T
b"T
product"T"
transpose_abool( "
transpose_bbool( "
Ttype:

2	
>
Maximum
x"T
y"T
z"T"
Ttype:
2	
?
Mean

input"T
reduction_indices"Tidx
output"T"
	keep_dimsbool( " 
Ttype:
2	"
Tidxtype0:
2	
e
MergeV2Checkpoints
checkpoint_prefixes
destination_prefix"
delete_old_dirsbool(?
>
Minimum
x"T
y"T
z"T"
Ttype:
2	
?
Mul
x"T
y"T
z"T"
Ttype:
2	?
0
Neg
x"T
y"T"
Ttype:
2
	

NoOp
M
Pack
values"T*N
output"T"
Nint(0"	
Ttype"
axisint 
C
Placeholder
output"dtype"
dtypetype"
shapeshape:
8
Pow
x"T
y"T
z"T"
Ttype:
2
	
N
PrintV2	
input"!
output_streamstringstderr"
endstring
?
e
Range
start"Tidx
limit"Tidx
delta"Tidx
output"Tidx"
Tidxtype0:
2		
@
ReadVariableOp
resource
value"dtype"
dtypetype?
@
RealDiv
x"T
y"T
z"T"
Ttype:
2	
7

Reciprocal
x"T
y"T"
Ttype:
2
	
E
Relu
features"T
activations"T"
Ttype:
2	
V
ReluGrad
	gradients"T
features"T
	backprops"T"
Ttype:
2	
[
Reshape
tensor"T
shape"Tshape
output"T"	
Ttype"
Tshapetype0:
2	
?
ResourceApplyAdam
var
m
v
beta1_power"T
beta2_power"T
lr"T

beta1"T

beta2"T
epsilon"T	
grad"T" 
Ttype:
2	"
use_lockingbool( "
use_nesterovbool( ?
o
	RestoreV2

prefix
tensor_names
shape_and_slices
tensors2dtypes"
dtypes
list(type)(0?
l
SaveV2

prefix
tensor_names
shape_and_slices
tensors2dtypes"
dtypes
list(type)(0?
?
Select
	condition

t"T
e"T
output"T"	
Ttype
A
SelectV2
	condition

t"T
e"T
output"T"	
Ttype
P
Shape

input"T
output"out_type"	
Ttype"
out_typetype0:
2	
H
ShardedFilename
basename	
shard

num_shards
filename
0
Sigmoid
x"T
y"T"
Ttype:

2
=
SigmoidGrad
y"T
dy"T
z"T"
Ttype:

2
-
Sqrt
x"T
y"T"
Ttype:

2
N
Squeeze

input"T
output"T"	
Ttype"
squeeze_dims	list(int)
 (
?
StatefulPartitionedCall
args2Tin
output2Tout"
Tin
list(type)("
Tout
list(type)("	
ffunc"
configstring "
config_protostring "
executor_typestring ??
@
StaticRegexFullMatch	
input

output
"
patternstring
?
StringFormat
inputs2T

output"
T
list(type)("
templatestring%s"
placeholderstring%s"
	summarizeint
N

StringJoin
inputs*N

output"
Nint(0"
	separatorstring 
<
Sub
x"T
y"T
z"T"
Ttype:
2	
?
Sum

input"T
reduction_indices"Tidx
output"T"
	keep_dimsbool( " 
Ttype:
2	"
Tidxtype0:
2	
c
Tile

input"T
	multiples"
Tmultiples
output"T"	
Ttype"

Tmultiplestype0:
2	
?
VarHandleOp
resource"
	containerstring "
shared_namestring "
dtypetype"
shapeshape"#
allowed_deviceslist(string)
 ?"serve*2.7.02v2.7.0-rc1-69-gc256c071bb28??
d
VariableVarHandleOp*
_output_shapes
: *
dtype0*
shape: *
shared_name
Variable
]
Variable/Read/ReadVariableOpReadVariableOpVariable*
_output_shapes
: *
dtype0
f
	Adam/iterVarHandleOp*
_output_shapes
: *
dtype0	*
shape: *
shared_name	Adam/iter
_
Adam/iter/Read/ReadVariableOpReadVariableOp	Adam/iter*
_output_shapes
: *
dtype0	
j
Adam/beta_1VarHandleOp*
_output_shapes
: *
dtype0*
shape: *
shared_nameAdam/beta_1
c
Adam/beta_1/Read/ReadVariableOpReadVariableOpAdam/beta_1*
_output_shapes
: *
dtype0
j
Adam/beta_2VarHandleOp*
_output_shapes
: *
dtype0*
shape: *
shared_nameAdam/beta_2
c
Adam/beta_2/Read/ReadVariableOpReadVariableOpAdam/beta_2*
_output_shapes
: *
dtype0
h

Adam/decayVarHandleOp*
_output_shapes
: *
dtype0*
shape: *
shared_name
Adam/decay
a
Adam/decay/Read/ReadVariableOpReadVariableOp
Adam/decay*
_output_shapes
: *
dtype0
x
Adam/learning_rateVarHandleOp*
_output_shapes
: *
dtype0*
shape: *#
shared_nameAdam/learning_rate
q
&Adam/learning_rate/Read/ReadVariableOpReadVariableOpAdam/learning_rate*
_output_shapes
: *
dtype0
t
dense/kernelVarHandleOp*
_output_shapes
: *
dtype0*
shape
:*
shared_namedense/kernel
m
 dense/kernel/Read/ReadVariableOpReadVariableOpdense/kernel*
_output_shapes

:*
dtype0
l

dense/biasVarHandleOp*
_output_shapes
: *
dtype0*
shape:*
shared_name
dense/bias
e
dense/bias/Read/ReadVariableOpReadVariableOp
dense/bias*
_output_shapes
:*
dtype0
x
dense_1/kernelVarHandleOp*
_output_shapes
: *
dtype0*
shape
:*
shared_namedense_1/kernel
q
"dense_1/kernel/Read/ReadVariableOpReadVariableOpdense_1/kernel*
_output_shapes

:*
dtype0
p
dense_1/biasVarHandleOp*
_output_shapes
: *
dtype0*
shape:*
shared_namedense_1/bias
i
 dense_1/bias/Read/ReadVariableOpReadVariableOpdense_1/bias*
_output_shapes
:*
dtype0
?
Adam/dense/kernel/mVarHandleOp*
_output_shapes
: *
dtype0*
shape
:*$
shared_nameAdam/dense/kernel/m
{
'Adam/dense/kernel/m/Read/ReadVariableOpReadVariableOpAdam/dense/kernel/m*
_output_shapes

:*
dtype0
z
Adam/dense/bias/mVarHandleOp*
_output_shapes
: *
dtype0*
shape:*"
shared_nameAdam/dense/bias/m
s
%Adam/dense/bias/m/Read/ReadVariableOpReadVariableOpAdam/dense/bias/m*
_output_shapes
:*
dtype0
?
Adam/dense_1/kernel/mVarHandleOp*
_output_shapes
: *
dtype0*
shape
:*&
shared_nameAdam/dense_1/kernel/m

)Adam/dense_1/kernel/m/Read/ReadVariableOpReadVariableOpAdam/dense_1/kernel/m*
_output_shapes

:*
dtype0
~
Adam/dense_1/bias/mVarHandleOp*
_output_shapes
: *
dtype0*
shape:*$
shared_nameAdam/dense_1/bias/m
w
'Adam/dense_1/bias/m/Read/ReadVariableOpReadVariableOpAdam/dense_1/bias/m*
_output_shapes
:*
dtype0
?
Adam/dense/kernel/vVarHandleOp*
_output_shapes
: *
dtype0*
shape
:*$
shared_nameAdam/dense/kernel/v
{
'Adam/dense/kernel/v/Read/ReadVariableOpReadVariableOpAdam/dense/kernel/v*
_output_shapes

:*
dtype0
z
Adam/dense/bias/vVarHandleOp*
_output_shapes
: *
dtype0*
shape:*"
shared_nameAdam/dense/bias/v
s
%Adam/dense/bias/v/Read/ReadVariableOpReadVariableOpAdam/dense/bias/v*
_output_shapes
:*
dtype0
?
Adam/dense_1/kernel/vVarHandleOp*
_output_shapes
: *
dtype0*
shape
:*&
shared_nameAdam/dense_1/kernel/v

)Adam/dense_1/kernel/v/Read/ReadVariableOpReadVariableOpAdam/dense_1/kernel/v*
_output_shapes

:*
dtype0
~
Adam/dense_1/bias/vVarHandleOp*
_output_shapes
: *
dtype0*
shape:*$
shared_nameAdam/dense_1/bias/v
w
'Adam/dense_1/bias/v/Read/ReadVariableOpReadVariableOpAdam/dense_1/bias/v*
_output_shapes
:*
dtype0

NoOpNoOp
?
ConstConst"/device:CPU:0*
_output_shapes
: *
dtype0*?
value?B? B?
>

_model
_global_step

_optimizer

signatures
?
layer_with_weights-0
layer-0
layer_with_weights-1
layer-1
	variables
trainable_variables
	regularization_losses

	keras_api
EC
VARIABLE_VALUEVariable'_global_step/.ATTRIBUTES/VARIABLE_VALUE
?
iter

beta_1

beta_2
	decay
learning_ratem+m,m-m.v/v0v1v2
 
h

kernel
bias
	variables
trainable_variables
regularization_losses
	keras_api
h

kernel
bias
	variables
trainable_variables
regularization_losses
	keras_api

0
1
2
3

0
1
2
3
 
?
non_trainable_variables

layers
metrics
layer_regularization_losses
 layer_metrics
	variables
trainable_variables
	regularization_losses
IG
VARIABLE_VALUE	Adam/iter*_optimizer/iter/.ATTRIBUTES/VARIABLE_VALUE
MK
VARIABLE_VALUEAdam/beta_1,_optimizer/beta_1/.ATTRIBUTES/VARIABLE_VALUE
MK
VARIABLE_VALUEAdam/beta_2,_optimizer/beta_2/.ATTRIBUTES/VARIABLE_VALUE
KI
VARIABLE_VALUE
Adam/decay+_optimizer/decay/.ATTRIBUTES/VARIABLE_VALUE
[Y
VARIABLE_VALUEAdam/learning_rate3_optimizer/learning_rate/.ATTRIBUTES/VARIABLE_VALUE
_]
VARIABLE_VALUEdense/kernel=_model/layer_with_weights-0/kernel/.ATTRIBUTES/VARIABLE_VALUE
[Y
VARIABLE_VALUE
dense/bias;_model/layer_with_weights-0/bias/.ATTRIBUTES/VARIABLE_VALUE

0
1

0
1
 
?
!non_trainable_variables

"layers
#metrics
$layer_regularization_losses
%layer_metrics
	variables
trainable_variables
regularization_losses
a_
VARIABLE_VALUEdense_1/kernel=_model/layer_with_weights-1/kernel/.ATTRIBUTES/VARIABLE_VALUE
][
VARIABLE_VALUEdense_1/bias;_model/layer_with_weights-1/bias/.ATTRIBUTES/VARIABLE_VALUE

0
1

0
1
 
?
&non_trainable_variables

'layers
(metrics
)layer_regularization_losses
*layer_metrics
	variables
trainable_variables
regularization_losses
 

0
1
 
 
 
 
 
 
 
 
 
 
 
 
 
??
VARIABLE_VALUEAdam/dense/kernel/mZ_model/layer_with_weights-0/kernel/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUE
}
VARIABLE_VALUEAdam/dense/bias/mX_model/layer_with_weights-0/bias/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUE
??
VARIABLE_VALUEAdam/dense_1/kernel/mZ_model/layer_with_weights-1/kernel/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUE
?
VARIABLE_VALUEAdam/dense_1/bias/mX_model/layer_with_weights-1/bias/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUE
??
VARIABLE_VALUEAdam/dense/kernel/vZ_model/layer_with_weights-0/kernel/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUE
}
VARIABLE_VALUEAdam/dense/bias/vX_model/layer_with_weights-0/bias/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUE
??
VARIABLE_VALUEAdam/dense_1/kernel/vZ_model/layer_with_weights-1/kernel/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUE
?
VARIABLE_VALUEAdam/dense_1/bias/vX_model/layer_with_weights-1/bias/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUE
m

learn_dataPlaceholder*'
_output_shapes
:?????????*
dtype0*
shape:?????????
g
learn_labelsPlaceholder*#
_output_shapes
:?????????*
dtype0*
shape:?????????
?
StatefulPartitionedCallStatefulPartitionedCall
learn_datalearn_labelsVariabledense/kernel
dense/biasdense_1/kerneldense_1/biasAdam/learning_rate	Adam/iterAdam/beta_1Adam/beta_2Adam/dense/kernel/mAdam/dense/kernel/vAdam/dense/bias/mAdam/dense/bias/vAdam/dense_1/kernel/mAdam/dense_1/kernel/vAdam/dense_1/bias/mAdam/dense_1/bias/v*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes
: *%
_read_only_resource_inputs
	
*-
config_proto

CPU

GPU 2J 8? **
f%R#
!__inference_signature_wrapper_768
o
predict_dataPlaceholder*'
_output_shapes
:?????????*
dtype0*
shape:?????????
?
StatefulPartitionedCall_1StatefulPartitionedCallpredict_datadense/kernel
dense/biasdense_1/kerneldense_1/bias*
Tin	
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *&
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? **
f%R#
!__inference_signature_wrapper_801
O
saver_filenamePlaceholder*
_output_shapes
: *
dtype0*
shape: 
?
StatefulPartitionedCall_2StatefulPartitionedCallsaver_filenameVariable/Read/ReadVariableOpAdam/iter/Read/ReadVariableOpAdam/beta_1/Read/ReadVariableOpAdam/beta_2/Read/ReadVariableOpAdam/decay/Read/ReadVariableOp&Adam/learning_rate/Read/ReadVariableOp dense/kernel/Read/ReadVariableOpdense/bias/Read/ReadVariableOp"dense_1/kernel/Read/ReadVariableOp dense_1/bias/Read/ReadVariableOp'Adam/dense/kernel/m/Read/ReadVariableOp%Adam/dense/bias/m/Read/ReadVariableOp)Adam/dense_1/kernel/m/Read/ReadVariableOp'Adam/dense_1/bias/m/Read/ReadVariableOp'Adam/dense/kernel/v/Read/ReadVariableOp%Adam/dense/bias/v/Read/ReadVariableOp)Adam/dense_1/kernel/v/Read/ReadVariableOp'Adam/dense_1/bias/v/Read/ReadVariableOpConst*
Tin
2	*
Tout
2*
_collective_manager_ids
 *
_output_shapes
: * 
_read_only_resource_inputs
 *-
config_proto

CPU

GPU 2J 8? *&
f!R
__inference__traced_save_1184
?
StatefulPartitionedCall_3StatefulPartitionedCallsaver_filenameVariable	Adam/iterAdam/beta_1Adam/beta_2
Adam/decayAdam/learning_ratedense/kernel
dense/biasdense_1/kerneldense_1/biasAdam/dense/kernel/mAdam/dense/bias/mAdam/dense_1/kernel/mAdam/dense_1/bias/mAdam/dense/kernel/vAdam/dense/bias/vAdam/dense_1/kernel/vAdam/dense_1/bias/v*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes
: * 
_read_only_resource_inputs
 *-
config_proto

CPU

GPU 2J 8? *)
f$R"
 __inference__traced_restore_1248??
?
?
D__inference_sequential_layer_call_and_return_conditional_losses_1036

inputs6
$dense_matmul_readvariableop_resource:3
%dense_biasadd_readvariableop_resource:8
&dense_1_matmul_readvariableop_resource:5
'dense_1_biasadd_readvariableop_resource:
identity??dense/BiasAdd/ReadVariableOp?dense/MatMul/ReadVariableOp?dense_1/BiasAdd/ReadVariableOp?dense_1/MatMul/ReadVariableOp?
dense/MatMul/ReadVariableOpReadVariableOp$dense_matmul_readvariableop_resource*
_output_shapes

:*
dtype0l
dense/MatMulMatMulinputs#dense/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ~
dense/BiasAdd/ReadVariableOpReadVariableOp%dense_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0
dense/BiasAddBiasAdddense/MatMul:product:0$dense/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: S

dense/ReluReludense/BiasAdd:output:0*
T0*
_output_shapes

: ?
dense_1/MatMul/ReadVariableOpReadVariableOp&dense_1_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
dense_1/MatMulMatMuldense/Relu:activations:0%dense_1/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
dense_1/BiasAdd/ReadVariableOpReadVariableOp'dense_1_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
dense_1/BiasAddBiasAdddense_1/MatMul:product:0&dense_1/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: ]
dense_1/SigmoidSigmoiddense_1/BiasAdd:output:0*
T0*
_output_shapes

: Y
IdentityIdentitydense_1/Sigmoid:y:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp^dense/BiasAdd/ReadVariableOp^dense/MatMul/ReadVariableOp^dense_1/BiasAdd/ReadVariableOp^dense_1/MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2<
dense/BiasAdd/ReadVariableOpdense/BiasAdd/ReadVariableOp2:
dense/MatMul/ReadVariableOpdense/MatMul/ReadVariableOp2@
dense_1/BiasAdd/ReadVariableOpdense_1/BiasAdd/ReadVariableOp2>
dense_1/MatMul/ReadVariableOpdense_1/MatMul/ReadVariableOp:F B

_output_shapes

: 
 
_user_specified_nameinputs
?	
?
A__inference_dense_1_layer_call_and_return_conditional_losses_1076

inputs0
matmul_readvariableop_resource:-
biasadd_readvariableop_resource:
identity??BiasAdd/ReadVariableOp?MatMul/ReadVariableOpt
MatMul/ReadVariableOpReadVariableOpmatmul_readvariableop_resource*
_output_shapes

:*
dtype0`
MatMulMatMulinputsMatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: r
BiasAdd/ReadVariableOpReadVariableOpbiasadd_readvariableop_resource*
_output_shapes
:*
dtype0m
BiasAddBiasAddMatMul:product:0BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: M
SigmoidSigmoidBiasAdd:output:0*
T0*
_output_shapes

: Q
IdentityIdentitySigmoid:y:0^NoOp*
T0*
_output_shapes

: w
NoOpNoOp^BiasAdd/ReadVariableOp^MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*!
_input_shapes
: : : 20
BiasAdd/ReadVariableOpBiasAdd/ReadVariableOp2.
MatMul/ReadVariableOpMatMul/ReadVariableOp:F B

_output_shapes

: 
 
_user_specified_nameinputs
?
?
__inference__wrapped_model_820
input_1A
/sequential_dense_matmul_readvariableop_resource:>
0sequential_dense_biasadd_readvariableop_resource:C
1sequential_dense_1_matmul_readvariableop_resource:@
2sequential_dense_1_biasadd_readvariableop_resource:
identity??'sequential/dense/BiasAdd/ReadVariableOp?&sequential/dense/MatMul/ReadVariableOp?)sequential/dense_1/BiasAdd/ReadVariableOp?(sequential/dense_1/MatMul/ReadVariableOp?
&sequential/dense/MatMul/ReadVariableOpReadVariableOp/sequential_dense_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
sequential/dense/MatMulMatMulinput_1.sequential/dense/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
'sequential/dense/BiasAdd/ReadVariableOpReadVariableOp0sequential_dense_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
sequential/dense/BiasAddBiasAdd!sequential/dense/MatMul:product:0/sequential/dense/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: i
sequential/dense/ReluRelu!sequential/dense/BiasAdd:output:0*
T0*
_output_shapes

: ?
(sequential/dense_1/MatMul/ReadVariableOpReadVariableOp1sequential_dense_1_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
sequential/dense_1/MatMulMatMul#sequential/dense/Relu:activations:00sequential/dense_1/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
)sequential/dense_1/BiasAdd/ReadVariableOpReadVariableOp2sequential_dense_1_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
sequential/dense_1/BiasAddBiasAdd#sequential/dense_1/MatMul:product:01sequential/dense_1/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: s
sequential/dense_1/SigmoidSigmoid#sequential/dense_1/BiasAdd:output:0*
T0*
_output_shapes

: d
IdentityIdentitysequential/dense_1/Sigmoid:y:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp(^sequential/dense/BiasAdd/ReadVariableOp'^sequential/dense/MatMul/ReadVariableOp*^sequential/dense_1/BiasAdd/ReadVariableOp)^sequential/dense_1/MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2R
'sequential/dense/BiasAdd/ReadVariableOp'sequential/dense/BiasAdd/ReadVariableOp2P
&sequential/dense/MatMul/ReadVariableOp&sequential/dense/MatMul/ReadVariableOp2V
)sequential/dense_1/BiasAdd/ReadVariableOp)sequential/dense_1/BiasAdd/ReadVariableOp2T
(sequential/dense_1/MatMul/ReadVariableOp(sequential/dense_1/MatMul/ReadVariableOp:G C

_output_shapes

: 
!
_user_specified_name	input_1
?
?
&__inference_dense_1_layer_call_fn_1065

inputs
unknown:
	unknown_0:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCallinputsunknown	unknown_0*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *I
fDRB
@__inference_dense_1_layer_call_and_return_conditional_losses_855f
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*!
_input_shapes
: : : 22
StatefulPartitionedCallStatefulPartitionedCall:F B

_output_shapes

: 
 
_user_specified_nameinputs
?-
?
__inference__traced_save_1184
file_prefix'
#savev2_variable_read_readvariableop(
$savev2_adam_iter_read_readvariableop	*
&savev2_adam_beta_1_read_readvariableop*
&savev2_adam_beta_2_read_readvariableop)
%savev2_adam_decay_read_readvariableop1
-savev2_adam_learning_rate_read_readvariableop+
'savev2_dense_kernel_read_readvariableop)
%savev2_dense_bias_read_readvariableop-
)savev2_dense_1_kernel_read_readvariableop+
'savev2_dense_1_bias_read_readvariableop2
.savev2_adam_dense_kernel_m_read_readvariableop0
,savev2_adam_dense_bias_m_read_readvariableop4
0savev2_adam_dense_1_kernel_m_read_readvariableop2
.savev2_adam_dense_1_bias_m_read_readvariableop2
.savev2_adam_dense_kernel_v_read_readvariableop0
,savev2_adam_dense_bias_v_read_readvariableop4
0savev2_adam_dense_1_kernel_v_read_readvariableop2
.savev2_adam_dense_1_bias_v_read_readvariableop
savev2_const

identity_1??MergeV2Checkpointsw
StaticRegexFullMatchStaticRegexFullMatchfile_prefix"/device:CPU:**
_output_shapes
: *
pattern
^s3://.*Z
ConstConst"/device:CPU:**
_output_shapes
: *
dtype0*
valueB B.parta
Const_1Const"/device:CPU:**
_output_shapes
: *
dtype0*
valueB B
_temp/part?
SelectSelectStaticRegexFullMatch:output:0Const:output:0Const_1:output:0"/device:CPU:**
T0*
_output_shapes
: f

StringJoin
StringJoinfile_prefixSelect:output:0"/device:CPU:**
N*
_output_shapes
: L

num_shardsConst*
_output_shapes
: *
dtype0*
value	B :f
ShardedFilename/shardConst"/device:CPU:0*
_output_shapes
: *
dtype0*
value	B : ?
ShardedFilenameShardedFilenameStringJoin:output:0ShardedFilename/shard:output:0num_shards:output:0"/device:CPU:0*
_output_shapes
: ?

SaveV2/tensor_namesConst"/device:CPU:0*
_output_shapes
:*
dtype0*?

value?
B?
B'_global_step/.ATTRIBUTES/VARIABLE_VALUEB*_optimizer/iter/.ATTRIBUTES/VARIABLE_VALUEB,_optimizer/beta_1/.ATTRIBUTES/VARIABLE_VALUEB,_optimizer/beta_2/.ATTRIBUTES/VARIABLE_VALUEB+_optimizer/decay/.ATTRIBUTES/VARIABLE_VALUEB3_optimizer/learning_rate/.ATTRIBUTES/VARIABLE_VALUEB=_model/layer_with_weights-0/kernel/.ATTRIBUTES/VARIABLE_VALUEB;_model/layer_with_weights-0/bias/.ATTRIBUTES/VARIABLE_VALUEB=_model/layer_with_weights-1/kernel/.ATTRIBUTES/VARIABLE_VALUEB;_model/layer_with_weights-1/bias/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-0/kernel/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-0/bias/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-1/kernel/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-1/bias/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-0/kernel/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-0/bias/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-1/kernel/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-1/bias/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEB_CHECKPOINTABLE_OBJECT_GRAPH?
SaveV2/shape_and_slicesConst"/device:CPU:0*
_output_shapes
:*
dtype0*9
value0B.B B B B B B B B B B B B B B B B B B B ?
SaveV2SaveV2ShardedFilename:filename:0SaveV2/tensor_names:output:0 SaveV2/shape_and_slices:output:0#savev2_variable_read_readvariableop$savev2_adam_iter_read_readvariableop&savev2_adam_beta_1_read_readvariableop&savev2_adam_beta_2_read_readvariableop%savev2_adam_decay_read_readvariableop-savev2_adam_learning_rate_read_readvariableop'savev2_dense_kernel_read_readvariableop%savev2_dense_bias_read_readvariableop)savev2_dense_1_kernel_read_readvariableop'savev2_dense_1_bias_read_readvariableop.savev2_adam_dense_kernel_m_read_readvariableop,savev2_adam_dense_bias_m_read_readvariableop0savev2_adam_dense_1_kernel_m_read_readvariableop.savev2_adam_dense_1_bias_m_read_readvariableop.savev2_adam_dense_kernel_v_read_readvariableop,savev2_adam_dense_bias_v_read_readvariableop0savev2_adam_dense_1_kernel_v_read_readvariableop.savev2_adam_dense_1_bias_v_read_readvariableopsavev2_const"/device:CPU:0*
_output_shapes
 *!
dtypes
2	?
&MergeV2Checkpoints/checkpoint_prefixesPackShardedFilename:filename:0^SaveV2"/device:CPU:0*
N*
T0*
_output_shapes
:?
MergeV2CheckpointsMergeV2Checkpoints/MergeV2Checkpoints/checkpoint_prefixes:output:0file_prefix"/device:CPU:0*
_output_shapes
 f
IdentityIdentityfile_prefix^MergeV2Checkpoints"/device:CPU:0*
T0*
_output_shapes
: Q

Identity_1IdentityIdentity:output:0^NoOp*
T0*
_output_shapes
: [
NoOpNoOp^MergeV2Checkpoints*"
_acd_function_control_output(*
_output_shapes
 "!

identity_1Identity_1:output:0*?
_input_shapesr
p: : : : : : : ::::::::::::: 2(
MergeV2CheckpointsMergeV2Checkpoints:C ?

_output_shapes
: 
%
_user_specified_namefile_prefix:

_output_shapes
: :

_output_shapes
: :

_output_shapes
: :

_output_shapes
: :

_output_shapes
: :

_output_shapes
: :$ 

_output_shapes

:: 

_output_shapes
::$	 

_output_shapes

:: 


_output_shapes
::$ 

_output_shapes

:: 

_output_shapes
::$ 

_output_shapes

:: 

_output_shapes
::$ 

_output_shapes

:: 

_output_shapes
::$ 

_output_shapes

:: 

_output_shapes
::

_output_shapes
: 
?
?
)__inference_sequential_layer_call_fn_1000

inputs
unknown:
	unknown_0:
	unknown_1:
	unknown_2:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCallinputsunknown	unknown_0	unknown_1	unknown_2*
Tin	
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *&
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *L
fGRE
C__inference_sequential_layer_call_and_return_conditional_losses_922f
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 22
StatefulPartitionedCallStatefulPartitionedCall:F B

_output_shapes

: 
 
_user_specified_nameinputs
?
?
C__inference_sequential_layer_call_and_return_conditional_losses_974
input_1
	dense_963:
	dense_965:
dense_1_968:
dense_1_970:
identity??dense/StatefulPartitionedCall?dense_1/StatefulPartitionedCall?
dense/StatefulPartitionedCallStatefulPartitionedCallinput_1	dense_963	dense_965*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *G
fBR@
>__inference_dense_layer_call_and_return_conditional_losses_838?
dense_1/StatefulPartitionedCallStatefulPartitionedCall&dense/StatefulPartitionedCall:output:0dense_1_968dense_1_970*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *I
fDRB
@__inference_dense_1_layer_call_and_return_conditional_losses_855n
IdentityIdentity(dense_1/StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp^dense/StatefulPartitionedCall ^dense_1/StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2>
dense/StatefulPartitionedCalldense/StatefulPartitionedCall2B
dense_1/StatefulPartitionedCalldense_1/StatefulPartitionedCall:G C

_output_shapes

: 
!
_user_specified_name	input_1
?
?
C__inference_sequential_layer_call_and_return_conditional_losses_862

inputs
	dense_839:
	dense_841:
dense_1_856:
dense_1_858:
identity??dense/StatefulPartitionedCall?dense_1/StatefulPartitionedCall?
dense/StatefulPartitionedCallStatefulPartitionedCallinputs	dense_839	dense_841*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *G
fBR@
>__inference_dense_layer_call_and_return_conditional_losses_838?
dense_1/StatefulPartitionedCallStatefulPartitionedCall&dense/StatefulPartitionedCall:output:0dense_1_856dense_1_858*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *I
fDRB
@__inference_dense_1_layer_call_and_return_conditional_losses_855n
IdentityIdentity(dense_1/StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp^dense/StatefulPartitionedCall ^dense_1/StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2>
dense/StatefulPartitionedCalldense/StatefulPartitionedCall2B
dense_1/StatefulPartitionedCalldense_1/StatefulPartitionedCall:F B

_output_shapes

: 
 
_user_specified_nameinputs
??
?
__inference_learn_726
data

labels&
assignaddvariableop_resource: A
/sequential_dense_matmul_readvariableop_resource:>
0sequential_dense_biasadd_readvariableop_resource:C
1sequential_dense_1_matmul_readvariableop_resource:@
2sequential_dense_1_biasadd_readvariableop_resource:+
!adam_cast_readvariableop_resource: &
adam_readvariableop_resource:	 -
#adam_cast_2_readvariableop_resource: -
#adam_cast_3_readvariableop_resource: 6
$adam_adam_update_resourceapplyadam_m:6
$adam_adam_update_resourceapplyadam_v:4
&adam_adam_update_1_resourceapplyadam_m:4
&adam_adam_update_1_resourceapplyadam_v:8
&adam_adam_update_2_resourceapplyadam_m:8
&adam_adam_update_2_resourceapplyadam_v:4
&adam_adam_update_3_resourceapplyadam_m:4
&adam_adam_update_3_resourceapplyadam_v:
identity??Adam/Adam/AssignAddVariableOp?"Adam/Adam/update/ResourceApplyAdam?$Adam/Adam/update_1/ResourceApplyAdam?$Adam/Adam/update_2/ResourceApplyAdam?$Adam/Adam/update_3/ResourceApplyAdam?Adam/Cast/ReadVariableOp?Adam/Cast_2/ReadVariableOp?Adam/Cast_3/ReadVariableOp?Adam/ReadVariableOp?AssignAddVariableOp?PrintV2?StringFormat/ReadVariableOp?'sequential/dense/BiasAdd/ReadVariableOp?&sequential/dense/MatMul/ReadVariableOp?)sequential/dense_1/BiasAdd/ReadVariableOp?(sequential/dense_1/MatMul/ReadVariableOpG
ConstConst*
_output_shapes
: *
dtype0*
value	B :{
AssignAddVariableOpAssignAddVariableOpassignaddvariableop_resourceConst:output:0*
_output_shapes
 *
dtype0?
&sequential/dense/MatMul/ReadVariableOpReadVariableOp/sequential_dense_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
sequential/dense/MatMulMatMuldata.sequential/dense/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
'sequential/dense/BiasAdd/ReadVariableOpReadVariableOp0sequential_dense_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
sequential/dense/BiasAddBiasAdd!sequential/dense/MatMul:product:0/sequential/dense/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: i
sequential/dense/ReluRelu!sequential/dense/BiasAdd:output:0*
T0*
_output_shapes

: ?
(sequential/dense_1/MatMul/ReadVariableOpReadVariableOp1sequential_dense_1_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
sequential/dense_1/MatMulMatMul#sequential/dense/Relu:activations:00sequential/dense_1/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
)sequential/dense_1/BiasAdd/ReadVariableOpReadVariableOp2sequential_dense_1_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
sequential/dense_1/BiasAddBiasAdd#sequential/dense_1/MatMul:product:01sequential/dense_1/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: s
sequential/dense_1/SigmoidSigmoid#sequential/dense_1/BiasAdd:output:0*
T0*
_output_shapes

: ?
8binary_crossentropy/remove_squeezable_dimensions/SqueezeSqueezesequential/dense_1/Sigmoid:y:0*
T0*
_output_shapes
: *
squeeze_dims

?????????e
binary_crossentropy/CastCastlabels*

DstT0*

SrcT0*#
_output_shapes
:?????????^
binary_crossentropy/ConstConst*
_output_shapes
: *
dtype0*
valueB
 *    `
binary_crossentropy/Const_1Const*
_output_shapes
: *
dtype0*
valueB
 *???3^
binary_crossentropy/sub/xConst*
_output_shapes
: *
dtype0*
valueB
 *  ???
binary_crossentropy/subSub"binary_crossentropy/sub/x:output:0$binary_crossentropy/Const_1:output:0*
T0*
_output_shapes
: ?
)binary_crossentropy/clip_by_value/MinimumMinimumAbinary_crossentropy/remove_squeezable_dimensions/Squeeze:output:0binary_crossentropy/sub:z:0*
T0*
_output_shapes
: ?
!binary_crossentropy/clip_by_valueMaximum-binary_crossentropy/clip_by_value/Minimum:z:0$binary_crossentropy/Const_1:output:0*
T0*
_output_shapes
: ^
binary_crossentropy/add/yConst*
_output_shapes
: *
dtype0*
valueB
 *???3?
binary_crossentropy/addAddV2%binary_crossentropy/clip_by_value:z:0"binary_crossentropy/add/y:output:0*
T0*
_output_shapes
: `
binary_crossentropy/LogLogbinary_crossentropy/add:z:0*
T0*
_output_shapes
: ~
binary_crossentropy/mulMulbinary_crossentropy/Cast:y:0binary_crossentropy/Log:y:0*
T0*
_output_shapes
: `
binary_crossentropy/sub_1/xConst*
_output_shapes
: *
dtype0*
valueB
 *  ???
binary_crossentropy/sub_1Sub$binary_crossentropy/sub_1/x:output:0binary_crossentropy/Cast:y:0*
T0*#
_output_shapes
:?????????`
binary_crossentropy/sub_2/xConst*
_output_shapes
: *
dtype0*
valueB
 *  ???
binary_crossentropy/sub_2Sub$binary_crossentropy/sub_2/x:output:0%binary_crossentropy/clip_by_value:z:0*
T0*
_output_shapes
: `
binary_crossentropy/add_1/yConst*
_output_shapes
: *
dtype0*
valueB
 *???3?
binary_crossentropy/add_1AddV2binary_crossentropy/sub_2:z:0$binary_crossentropy/add_1/y:output:0*
T0*
_output_shapes
: d
binary_crossentropy/Log_1Logbinary_crossentropy/add_1:z:0*
T0*
_output_shapes
: ?
binary_crossentropy/mul_1Mulbinary_crossentropy/sub_1:z:0binary_crossentropy/Log_1:y:0*
T0*
_output_shapes
: ?
binary_crossentropy/add_2AddV2binary_crossentropy/mul:z:0binary_crossentropy/mul_1:z:0*
T0*
_output_shapes
: b
binary_crossentropy/NegNegbinary_crossentropy/add_2:z:0*
T0*
_output_shapes
: u
*binary_crossentropy/Mean/reduction_indicesConst*
_output_shapes
: *
dtype0*
valueB :
??????????
binary_crossentropy/MeanMeanbinary_crossentropy/Neg:y:03binary_crossentropy/Mean/reduction_indices:output:0*
T0*
_output_shapes
: l
'binary_crossentropy/weighted_loss/ConstConst*
_output_shapes
: *
dtype0*
valueB
 *  ???
%binary_crossentropy/weighted_loss/MulMul!binary_crossentropy/Mean:output:00binary_crossentropy/weighted_loss/Const:output:0*
T0*
_output_shapes
: h
&binary_crossentropy/weighted_loss/RankConst*
_output_shapes
: *
dtype0*
value	B : o
-binary_crossentropy/weighted_loss/range/startConst*
_output_shapes
: *
dtype0*
value	B : o
-binary_crossentropy/weighted_loss/range/deltaConst*
_output_shapes
: *
dtype0*
value	B :?
'binary_crossentropy/weighted_loss/rangeRange6binary_crossentropy/weighted_loss/range/start:output:0/binary_crossentropy/weighted_loss/Rank:output:06binary_crossentropy/weighted_loss/range/delta:output:0*
_output_shapes
: ?
%binary_crossentropy/weighted_loss/SumSum)binary_crossentropy/weighted_loss/Mul:z:00binary_crossentropy/weighted_loss/range:output:0*
T0*
_output_shapes
: p
.binary_crossentropy/weighted_loss/num_elementsConst*
_output_shapes
: *
dtype0*
value	B :?
3binary_crossentropy/weighted_loss/num_elements/CastCast7binary_crossentropy/weighted_loss/num_elements:output:0*

DstT0*

SrcT0*
_output_shapes
: j
(binary_crossentropy/weighted_loss/Rank_1Const*
_output_shapes
: *
dtype0*
value	B : q
/binary_crossentropy/weighted_loss/range_1/startConst*
_output_shapes
: *
dtype0*
value	B : q
/binary_crossentropy/weighted_loss/range_1/deltaConst*
_output_shapes
: *
dtype0*
value	B :?
)binary_crossentropy/weighted_loss/range_1Range8binary_crossentropy/weighted_loss/range_1/start:output:01binary_crossentropy/weighted_loss/Rank_1:output:08binary_crossentropy/weighted_loss/range_1/delta:output:0*
_output_shapes
: ?
'binary_crossentropy/weighted_loss/Sum_1Sum.binary_crossentropy/weighted_loss/Sum:output:02binary_crossentropy/weighted_loss/range_1:output:0*
T0*
_output_shapes
: ?
'binary_crossentropy/weighted_loss/valueDivNoNan0binary_crossentropy/weighted_loss/Sum_1:output:07binary_crossentropy/weighted_loss/num_elements/Cast:y:0*
T0*
_output_shapes
: ?
StringFormat/ReadVariableOpReadVariableOpassignaddvariableop_resource^AssignAddVariableOp*
_output_shapes
: *
dtype0?
StringFormatStringFormat#StringFormat/ReadVariableOp:value:0+binary_crossentropy/weighted_loss/value:z:0*
T
2*
_output_shapes
: *
placeholder{}*
template{} : loss:  {}?
PrintV2PrintV2StringFormat:output:0*
_output_shapes
 I
onesConst*
_output_shapes
: *
dtype0*
valueB
 *  ??~
;gradient_tape/binary_crossentropy/weighted_loss/value/ShapeConst*
_output_shapes
: *
dtype0*
valueB ?
=gradient_tape/binary_crossentropy/weighted_loss/value/Shape_1Const*
_output_shapes
: *
dtype0*
valueB ?
Kgradient_tape/binary_crossentropy/weighted_loss/value/BroadcastGradientArgsBroadcastGradientArgsDgradient_tape/binary_crossentropy/weighted_loss/value/Shape:output:0Fgradient_tape/binary_crossentropy/weighted_loss/value/Shape_1:output:0*2
_output_shapes 
:?????????:??????????
@gradient_tape/binary_crossentropy/weighted_loss/value/div_no_nanDivNoNanones:output:07binary_crossentropy/weighted_loss/num_elements/Cast:y:0*
T0*
_output_shapes
: ?
9gradient_tape/binary_crossentropy/weighted_loss/value/SumSumDgradient_tape/binary_crossentropy/weighted_loss/value/div_no_nan:z:0Pgradient_tape/binary_crossentropy/weighted_loss/value/BroadcastGradientArgs:r0:0*
T0*
_output_shapes
: ?
=gradient_tape/binary_crossentropy/weighted_loss/value/ReshapeReshapeBgradient_tape/binary_crossentropy/weighted_loss/value/Sum:output:0Dgradient_tape/binary_crossentropy/weighted_loss/value/Shape:output:0*
T0*
_output_shapes
: ?
9gradient_tape/binary_crossentropy/weighted_loss/value/NegNeg0binary_crossentropy/weighted_loss/Sum_1:output:0*
T0*
_output_shapes
: ?
Bgradient_tape/binary_crossentropy/weighted_loss/value/div_no_nan_1DivNoNan=gradient_tape/binary_crossentropy/weighted_loss/value/Neg:y:07binary_crossentropy/weighted_loss/num_elements/Cast:y:0*
T0*
_output_shapes
: ?
Bgradient_tape/binary_crossentropy/weighted_loss/value/div_no_nan_2DivNoNanFgradient_tape/binary_crossentropy/weighted_loss/value/div_no_nan_1:z:07binary_crossentropy/weighted_loss/num_elements/Cast:y:0*
T0*
_output_shapes
: ?
9gradient_tape/binary_crossentropy/weighted_loss/value/mulMulones:output:0Fgradient_tape/binary_crossentropy/weighted_loss/value/div_no_nan_2:z:0*
T0*
_output_shapes
: ?
;gradient_tape/binary_crossentropy/weighted_loss/value/Sum_1Sum=gradient_tape/binary_crossentropy/weighted_loss/value/mul:z:0Pgradient_tape/binary_crossentropy/weighted_loss/value/BroadcastGradientArgs:r1:0*
T0*
_output_shapes
: ?
?gradient_tape/binary_crossentropy/weighted_loss/value/Reshape_1ReshapeDgradient_tape/binary_crossentropy/weighted_loss/value/Sum_1:output:0Fgradient_tape/binary_crossentropy/weighted_loss/value/Shape_1:output:0*
T0*
_output_shapes
: ?
=gradient_tape/binary_crossentropy/weighted_loss/Reshape/shapeConst*
_output_shapes
: *
dtype0*
valueB ?
?gradient_tape/binary_crossentropy/weighted_loss/Reshape/shape_1Const*
_output_shapes
: *
dtype0*
valueB ?
7gradient_tape/binary_crossentropy/weighted_loss/ReshapeReshapeFgradient_tape/binary_crossentropy/weighted_loss/value/Reshape:output:0Hgradient_tape/binary_crossentropy/weighted_loss/Reshape/shape_1:output:0*
T0*
_output_shapes
: x
5gradient_tape/binary_crossentropy/weighted_loss/ConstConst*
_output_shapes
: *
dtype0*
valueB ?
4gradient_tape/binary_crossentropy/weighted_loss/TileTile@gradient_tape/binary_crossentropy/weighted_loss/Reshape:output:0>gradient_tape/binary_crossentropy/weighted_loss/Const:output:0*
T0*
_output_shapes
: ?
?gradient_tape/binary_crossentropy/weighted_loss/Reshape_1/shapeConst*
_output_shapes
: *
dtype0*
valueB ?
Agradient_tape/binary_crossentropy/weighted_loss/Reshape_1/shape_1Const*
_output_shapes
: *
dtype0*
valueB ?
9gradient_tape/binary_crossentropy/weighted_loss/Reshape_1Reshape=gradient_tape/binary_crossentropy/weighted_loss/Tile:output:0Jgradient_tape/binary_crossentropy/weighted_loss/Reshape_1/shape_1:output:0*
T0*
_output_shapes
: z
7gradient_tape/binary_crossentropy/weighted_loss/Const_1Const*
_output_shapes
: *
dtype0*
valueB ?
6gradient_tape/binary_crossentropy/weighted_loss/Tile_1TileBgradient_tape/binary_crossentropy/weighted_loss/Reshape_1:output:0@gradient_tape/binary_crossentropy/weighted_loss/Const_1:output:0*
T0*
_output_shapes
: ?
3gradient_tape/binary_crossentropy/weighted_loss/MulMul?gradient_tape/binary_crossentropy/weighted_loss/Tile_1:output:00binary_crossentropy/weighted_loss/Const:output:0*
T0*
_output_shapes
: u
+gradient_tape/binary_crossentropy/Maximum/xConst*
_output_shapes
:*
dtype0*
valueB:m
+gradient_tape/binary_crossentropy/Maximum/yConst*
_output_shapes
: *
dtype0*
value	B :?
)gradient_tape/binary_crossentropy/MaximumMaximum4gradient_tape/binary_crossentropy/Maximum/x:output:04gradient_tape/binary_crossentropy/Maximum/y:output:0*
T0*
_output_shapes
:v
,gradient_tape/binary_crossentropy/floordiv/xConst*
_output_shapes
:*
dtype0*
valueB: ?
*gradient_tape/binary_crossentropy/floordivFloorDiv5gradient_tape/binary_crossentropy/floordiv/x:output:0-gradient_tape/binary_crossentropy/Maximum:z:0*
T0*
_output_shapes
:y
/gradient_tape/binary_crossentropy/Reshape/shapeConst*
_output_shapes
:*
dtype0*
valueB:?
)gradient_tape/binary_crossentropy/ReshapeReshape7gradient_tape/binary_crossentropy/weighted_loss/Mul:z:08gradient_tape/binary_crossentropy/Reshape/shape:output:0*
T0*
_output_shapes
:z
0gradient_tape/binary_crossentropy/Tile/multiplesConst*
_output_shapes
:*
dtype0*
valueB: ?
&gradient_tape/binary_crossentropy/TileTile2gradient_tape/binary_crossentropy/Reshape:output:09gradient_tape/binary_crossentropy/Tile/multiples:output:0*
T0*
_output_shapes
: l
'gradient_tape/binary_crossentropy/ConstConst*
_output_shapes
: *
dtype0*
valueB
 *   B?
)gradient_tape/binary_crossentropy/truedivRealDiv/gradient_tape/binary_crossentropy/Tile:output:00gradient_tape/binary_crossentropy/Const:output:0*
T0*
_output_shapes
: ?
%gradient_tape/binary_crossentropy/NegNeg-gradient_tape/binary_crossentropy/truediv:z:0*
T0*
_output_shapes
: w
+gradient_tape/binary_crossentropy/mul/ShapeShapebinary_crossentropy/Cast:y:0*
T0*
_output_shapes
:x
-gradient_tape/binary_crossentropy/mul/Shape_1Shapebinary_crossentropy/Log:y:0*
T0*
_output_shapes
:?
;gradient_tape/binary_crossentropy/mul/BroadcastGradientArgsBroadcastGradientArgs4gradient_tape/binary_crossentropy/mul/Shape:output:06gradient_tape/binary_crossentropy/mul/Shape_1:output:0*2
_output_shapes 
:?????????:??????????
)gradient_tape/binary_crossentropy/mul/MulMulbinary_crossentropy/Cast:y:0)gradient_tape/binary_crossentropy/Neg:y:0*
T0*
_output_shapes
: ?
)gradient_tape/binary_crossentropy/mul/SumSum-gradient_tape/binary_crossentropy/mul/Mul:z:0@gradient_tape/binary_crossentropy/mul/BroadcastGradientArgs:r1:0*
T0*
_output_shapes
:?
-gradient_tape/binary_crossentropy/mul/ReshapeReshape2gradient_tape/binary_crossentropy/mul/Sum:output:06gradient_tape/binary_crossentropy/mul/Shape_1:output:0*
T0*
_output_shapes
: z
-gradient_tape/binary_crossentropy/mul_1/ShapeShapebinary_crossentropy/sub_1:z:0*
T0*
_output_shapes
:|
/gradient_tape/binary_crossentropy/mul_1/Shape_1Shapebinary_crossentropy/Log_1:y:0*
T0*
_output_shapes
:?
=gradient_tape/binary_crossentropy/mul_1/BroadcastGradientArgsBroadcastGradientArgs6gradient_tape/binary_crossentropy/mul_1/Shape:output:08gradient_tape/binary_crossentropy/mul_1/Shape_1:output:0*2
_output_shapes 
:?????????:??????????
+gradient_tape/binary_crossentropy/mul_1/MulMulbinary_crossentropy/sub_1:z:0)gradient_tape/binary_crossentropy/Neg:y:0*
T0*
_output_shapes
: ?
+gradient_tape/binary_crossentropy/mul_1/SumSum/gradient_tape/binary_crossentropy/mul_1/Mul:z:0Bgradient_tape/binary_crossentropy/mul_1/BroadcastGradientArgs:r1:0*
T0*
_output_shapes
:?
/gradient_tape/binary_crossentropy/mul_1/ReshapeReshape4gradient_tape/binary_crossentropy/mul_1/Sum:output:08gradient_tape/binary_crossentropy/mul_1/Shape_1:output:0*
T0*
_output_shapes
: ?
,gradient_tape/binary_crossentropy/Reciprocal
Reciprocalbinary_crossentropy/add:z:0.^gradient_tape/binary_crossentropy/mul/Reshape*
T0*
_output_shapes
: ?
%gradient_tape/binary_crossentropy/mulMul6gradient_tape/binary_crossentropy/mul/Reshape:output:00gradient_tape/binary_crossentropy/Reciprocal:y:0*
T0*
_output_shapes
: ?
.gradient_tape/binary_crossentropy/Reciprocal_1
Reciprocalbinary_crossentropy/add_1:z:00^gradient_tape/binary_crossentropy/mul_1/Reshape*
T0*
_output_shapes
: ?
'gradient_tape/binary_crossentropy/mul_1Mul8gradient_tape/binary_crossentropy/mul_1/Reshape:output:02gradient_tape/binary_crossentropy/Reciprocal_1:y:0*
T0*
_output_shapes
: ?
@gradient_tape/binary_crossentropy/sub_2/BroadcastGradientArgs/s0Const*
_output_shapes
: *
dtype0*
valueB ?
Bgradient_tape/binary_crossentropy/sub_2/BroadcastGradientArgs/s0_1Const*
_output_shapes
: *
dtype0*
valueB ?
@gradient_tape/binary_crossentropy/sub_2/BroadcastGradientArgs/s1Const*
_output_shapes
:*
dtype0*
valueB: ?
=gradient_tape/binary_crossentropy/sub_2/BroadcastGradientArgsBroadcastGradientArgsKgradient_tape/binary_crossentropy/sub_2/BroadcastGradientArgs/s0_1:output:0Igradient_tape/binary_crossentropy/sub_2/BroadcastGradientArgs/s1:output:0*2
_output_shapes 
:?????????:??????????
+gradient_tape/binary_crossentropy/sub_2/NegNeg+gradient_tape/binary_crossentropy/mul_1:z:0*
T0*
_output_shapes
: ?
AddNAddN)gradient_tape/binary_crossentropy/mul:z:0/gradient_tape/binary_crossentropy/sub_2/Neg:y:0*
N*
T0*
_output_shapes
: ?
:gradient_tape/binary_crossentropy/clip_by_value/zeros_likeConst*
_output_shapes
: *
dtype0*
valueB *    ?
<gradient_tape/binary_crossentropy/clip_by_value/GreaterEqualGreaterEqual-binary_crossentropy/clip_by_value/Minimum:z:0$binary_crossentropy/Const_1:output:0*
T0*
_output_shapes
: ?
8gradient_tape/binary_crossentropy/clip_by_value/SelectV2SelectV2@gradient_tape/binary_crossentropy/clip_by_value/GreaterEqual:z:0
AddN:sum:0Cgradient_tape/binary_crossentropy/clip_by_value/zeros_like:output:0*
T0*
_output_shapes
: ?
<gradient_tape/binary_crossentropy/clip_by_value/zeros_like_1Const*
_output_shapes
: *
dtype0*
valueB *    ?
9gradient_tape/binary_crossentropy/clip_by_value/LessEqual	LessEqualAbinary_crossentropy/remove_squeezable_dimensions/Squeeze:output:0binary_crossentropy/sub:z:0*
T0*
_output_shapes
: ?
:gradient_tape/binary_crossentropy/clip_by_value/SelectV2_1SelectV2=gradient_tape/binary_crossentropy/clip_by_value/LessEqual:z:0Agradient_tape/binary_crossentropy/clip_by_value/SelectV2:output:0Egradient_tape/binary_crossentropy/clip_by_value/zeros_like_1:output:0*
T0*
_output_shapes
: ?
Dgradient_tape/binary_crossentropy/remove_squeezable_dimensions/ShapeConst*
_output_shapes
:*
dtype0*
valueB"       ?
Fgradient_tape/binary_crossentropy/remove_squeezable_dimensions/ReshapeReshapeCgradient_tape/binary_crossentropy/clip_by_value/SelectV2_1:output:0Mgradient_tape/binary_crossentropy/remove_squeezable_dimensions/Shape:output:0*
T0*
_output_shapes

: ?
4gradient_tape/sequential/dense_1/Sigmoid/SigmoidGradSigmoidGradsequential/dense_1/Sigmoid:y:0Ogradient_tape/binary_crossentropy/remove_squeezable_dimensions/Reshape:output:0*
T0*
_output_shapes

: ?
4gradient_tape/sequential/dense_1/BiasAdd/BiasAddGradBiasAddGrad8gradient_tape/sequential/dense_1/Sigmoid/SigmoidGrad:z:0*
T0*
_output_shapes
:?
.gradient_tape/sequential/dense_1/MatMul/MatMulMatMul8gradient_tape/sequential/dense_1/Sigmoid/SigmoidGrad:z:00sequential/dense_1/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: *
transpose_b(?
0gradient_tape/sequential/dense_1/MatMul/MatMul_1MatMul#sequential/dense/Relu:activations:08gradient_tape/sequential/dense_1/Sigmoid/SigmoidGrad:z:0*
T0*
_output_shapes

:*
transpose_a(?
'gradient_tape/sequential/dense/ReluGradReluGrad8gradient_tape/sequential/dense_1/MatMul/MatMul:product:0#sequential/dense/Relu:activations:0*
T0*
_output_shapes

: ?
2gradient_tape/sequential/dense/BiasAdd/BiasAddGradBiasAddGrad3gradient_tape/sequential/dense/ReluGrad:backprops:0*
T0*
_output_shapes
:?
,gradient_tape/sequential/dense/MatMul/MatMulMatMuldata3gradient_tape/sequential/dense/ReluGrad:backprops:0*
T0*
_output_shapes

:*
transpose_a(r
Adam/Cast/ReadVariableOpReadVariableOp!adam_cast_readvariableop_resource*
_output_shapes
: *
dtype0?
Adam/IdentityIdentity Adam/Cast/ReadVariableOp:value:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: h
Adam/ReadVariableOpReadVariableOpadam_readvariableop_resource*
_output_shapes
: *
dtype0	z

Adam/add/yConst",/job:localhost/replica:0/task:0/device:CPU:0*
_output_shapes
: *
dtype0	*
value	B	 R?
Adam/addAddV2Adam/ReadVariableOp:value:0Adam/add/y:output:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0	*
_output_shapes
: 
Adam/Cast_1CastAdam/add:z:0",/job:localhost/replica:0/task:0/device:CPU:0*

DstT0*

SrcT0	*
_output_shapes
: v
Adam/Cast_2/ReadVariableOpReadVariableOp#adam_cast_2_readvariableop_resource*
_output_shapes
: *
dtype0?
Adam/Identity_1Identity"Adam/Cast_2/ReadVariableOp:value:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: v
Adam/Cast_3/ReadVariableOpReadVariableOp#adam_cast_3_readvariableop_resource*
_output_shapes
: *
dtype0?
Adam/Identity_2Identity"Adam/Cast_3/ReadVariableOp:value:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: ?
Adam/PowPowAdam/Identity_1:output:0Adam/Cast_1:y:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: ?

Adam/Pow_1PowAdam/Identity_2:output:0Adam/Cast_1:y:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: }

Adam/sub/xConst",/job:localhost/replica:0/task:0/device:CPU:0*
_output_shapes
: *
dtype0*
valueB
 *  ???
Adam/subSubAdam/sub/x:output:0Adam/Pow_1:z:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: n
	Adam/SqrtSqrtAdam/sub:z:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: 
Adam/sub_1/xConst",/job:localhost/replica:0/task:0/device:CPU:0*
_output_shapes
: *
dtype0*
valueB
 *  ???

Adam/sub_1SubAdam/sub_1/x:output:0Adam/Pow:z:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: ?
Adam/truedivRealDivAdam/Sqrt:y:0Adam/sub_1:z:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: ?
Adam/mulMulAdam/Identity:output:0Adam/truediv:z:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: }

Adam/ConstConst",/job:localhost/replica:0/task:0/device:CPU:0*
_output_shapes
: *
dtype0*
valueB
 *???3
Adam/sub_2/xConst",/job:localhost/replica:0/task:0/device:CPU:0*
_output_shapes
: *
dtype0*
valueB
 *  ???

Adam/sub_2SubAdam/sub_2/x:output:0Adam/Identity_1:output:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: 
Adam/sub_3/xConst",/job:localhost/replica:0/task:0/device:CPU:0*
_output_shapes
: *
dtype0*
valueB
 *  ???

Adam/sub_3SubAdam/sub_3/x:output:0Adam/Identity_2:output:0",/job:localhost/replica:0/task:0/device:CPU:0*
T0*
_output_shapes
: |
Adam/Identity_3Identity6gradient_tape/sequential/dense/MatMul/MatMul:product:0*
T0*
_output_shapes

:}
Adam/Identity_4Identity;gradient_tape/sequential/dense/BiasAdd/BiasAddGrad:output:0*
T0*
_output_shapes
:?
Adam/Identity_5Identity:gradient_tape/sequential/dense_1/MatMul/MatMul_1:product:0*
T0*
_output_shapes

:
Adam/Identity_6Identity=gradient_tape/sequential/dense_1/BiasAdd/BiasAddGrad:output:0*
T0*
_output_shapes
:?
Adam/IdentityN	IdentityN6gradient_tape/sequential/dense/MatMul/MatMul:product:0;gradient_tape/sequential/dense/BiasAdd/BiasAddGrad:output:0:gradient_tape/sequential/dense_1/MatMul/MatMul_1:product:0=gradient_tape/sequential/dense_1/BiasAdd/BiasAddGrad:output:06gradient_tape/sequential/dense/MatMul/MatMul:product:0;gradient_tape/sequential/dense/BiasAdd/BiasAddGrad:output:0:gradient_tape/sequential/dense_1/MatMul/MatMul_1:product:0=gradient_tape/sequential/dense_1/BiasAdd/BiasAddGrad:output:0*
T

2*)
_gradient_op_typeCustomGradient-699*T
_output_shapesB
@::::::::?
"Adam/Adam/update/ResourceApplyAdamResourceApplyAdam/sequential_dense_matmul_readvariableop_resource$adam_adam_update_resourceapplyadam_m$adam_adam_update_resourceapplyadam_vAdam/Pow:z:0Adam/Pow_1:z:0Adam/Identity:output:0Adam/Identity_1:output:0Adam/Identity_2:output:0Adam/Const:output:0Adam/IdentityN:output:0'^sequential/dense/MatMul/ReadVariableOp",/job:localhost/replica:0/task:0/device:CPU:0*
T0*B
_class8
64loc:@sequential/dense/MatMul/ReadVariableOp/resource*
_output_shapes
 *
use_locking(?
$Adam/Adam/update_1/ResourceApplyAdamResourceApplyAdam0sequential_dense_biasadd_readvariableop_resource&adam_adam_update_1_resourceapplyadam_m&adam_adam_update_1_resourceapplyadam_vAdam/Pow:z:0Adam/Pow_1:z:0Adam/Identity:output:0Adam/Identity_1:output:0Adam/Identity_2:output:0Adam/Const:output:0Adam/IdentityN:output:1(^sequential/dense/BiasAdd/ReadVariableOp",/job:localhost/replica:0/task:0/device:CPU:0*
T0*C
_class9
75loc:@sequential/dense/BiasAdd/ReadVariableOp/resource*
_output_shapes
 *
use_locking(?
$Adam/Adam/update_2/ResourceApplyAdamResourceApplyAdam1sequential_dense_1_matmul_readvariableop_resource&adam_adam_update_2_resourceapplyadam_m&adam_adam_update_2_resourceapplyadam_vAdam/Pow:z:0Adam/Pow_1:z:0Adam/Identity:output:0Adam/Identity_1:output:0Adam/Identity_2:output:0Adam/Const:output:0Adam/IdentityN:output:2)^sequential/dense_1/MatMul/ReadVariableOp",/job:localhost/replica:0/task:0/device:CPU:0*
T0*D
_class:
86loc:@sequential/dense_1/MatMul/ReadVariableOp/resource*
_output_shapes
 *
use_locking(?
$Adam/Adam/update_3/ResourceApplyAdamResourceApplyAdam2sequential_dense_1_biasadd_readvariableop_resource&adam_adam_update_3_resourceapplyadam_m&adam_adam_update_3_resourceapplyadam_vAdam/Pow:z:0Adam/Pow_1:z:0Adam/Identity:output:0Adam/Identity_1:output:0Adam/Identity_2:output:0Adam/Const:output:0Adam/IdentityN:output:3*^sequential/dense_1/BiasAdd/ReadVariableOp",/job:localhost/replica:0/task:0/device:CPU:0*
T0*E
_class;
97loc:@sequential/dense_1/BiasAdd/ReadVariableOp/resource*
_output_shapes
 *
use_locking(?
Adam/Adam/group_depsNoOp#^Adam/Adam/update/ResourceApplyAdam%^Adam/Adam/update_1/ResourceApplyAdam%^Adam/Adam/update_2/ResourceApplyAdam%^Adam/Adam/update_3/ResourceApplyAdam",/job:localhost/replica:0/task:0/device:CPU:0*
_output_shapes
 h
Adam/Adam/ConstConst^Adam/Adam/group_deps*
_output_shapes
: *
dtype0	*
value	B	 R?
Adam/Adam/AssignAddVariableOpAssignAddVariableOpadam_readvariableop_resourceAdam/Adam/Const:output:0^Adam/ReadVariableOp*
_output_shapes
 *
dtype0	i
IdentityIdentity+binary_crossentropy/weighted_loss/value:z:0^NoOp*
T0*
_output_shapes
: ?
NoOpNoOp^Adam/Adam/AssignAddVariableOp#^Adam/Adam/update/ResourceApplyAdam%^Adam/Adam/update_1/ResourceApplyAdam%^Adam/Adam/update_2/ResourceApplyAdam%^Adam/Adam/update_3/ResourceApplyAdam^Adam/Cast/ReadVariableOp^Adam/Cast_2/ReadVariableOp^Adam/Cast_3/ReadVariableOp^Adam/ReadVariableOp^AssignAddVariableOp^PrintV2^StringFormat/ReadVariableOp(^sequential/dense/BiasAdd/ReadVariableOp'^sequential/dense/MatMul/ReadVariableOp*^sequential/dense_1/BiasAdd/ReadVariableOp)^sequential/dense_1/MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*N
_input_shapes=
;: :?????????: : : : : : : : : : : : : : : : : 2>
Adam/Adam/AssignAddVariableOpAdam/Adam/AssignAddVariableOp2H
"Adam/Adam/update/ResourceApplyAdam"Adam/Adam/update/ResourceApplyAdam2L
$Adam/Adam/update_1/ResourceApplyAdam$Adam/Adam/update_1/ResourceApplyAdam2L
$Adam/Adam/update_2/ResourceApplyAdam$Adam/Adam/update_2/ResourceApplyAdam2L
$Adam/Adam/update_3/ResourceApplyAdam$Adam/Adam/update_3/ResourceApplyAdam24
Adam/Cast/ReadVariableOpAdam/Cast/ReadVariableOp28
Adam/Cast_2/ReadVariableOpAdam/Cast_2/ReadVariableOp28
Adam/Cast_3/ReadVariableOpAdam/Cast_3/ReadVariableOp2*
Adam/ReadVariableOpAdam/ReadVariableOp2*
AssignAddVariableOpAssignAddVariableOp2
PrintV2PrintV22:
StringFormat/ReadVariableOpStringFormat/ReadVariableOp2R
'sequential/dense/BiasAdd/ReadVariableOp'sequential/dense/BiasAdd/ReadVariableOp2P
&sequential/dense/MatMul/ReadVariableOp&sequential/dense/MatMul/ReadVariableOp2V
)sequential/dense_1/BiasAdd/ReadVariableOp)sequential/dense_1/BiasAdd/ReadVariableOp2T
(sequential/dense_1/MatMul/ReadVariableOp(sequential/dense_1/MatMul/ReadVariableOp:M I
'
_output_shapes
:?????????

_user_specified_namedata:KG
#
_output_shapes
:?????????
 
_user_specified_namelabels:HD
B
_class8
64loc:@sequential/dense/MatMul/ReadVariableOp/resource:HD
B
_class8
64loc:@sequential/dense/MatMul/ReadVariableOp/resource:IE
C
_class9
75loc:@sequential/dense/BiasAdd/ReadVariableOp/resource:IE
C
_class9
75loc:@sequential/dense/BiasAdd/ReadVariableOp/resource:JF
D
_class:
86loc:@sequential/dense_1/MatMul/ReadVariableOp/resource:JF
D
_class:
86loc:@sequential/dense_1/MatMul/ReadVariableOp/resource:KG
E
_class;
97loc:@sequential/dense_1/BiasAdd/ReadVariableOp/resource:KG
E
_class;
97loc:@sequential/dense_1/BiasAdd/ReadVariableOp/resource
?	
?
?__inference_dense_layer_call_and_return_conditional_losses_1056

inputs0
matmul_readvariableop_resource:-
biasadd_readvariableop_resource:
identity??BiasAdd/ReadVariableOp?MatMul/ReadVariableOpt
MatMul/ReadVariableOpReadVariableOpmatmul_readvariableop_resource*
_output_shapes

:*
dtype0`
MatMulMatMulinputsMatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: r
BiasAdd/ReadVariableOpReadVariableOpbiasadd_readvariableop_resource*
_output_shapes
:*
dtype0m
BiasAddBiasAddMatMul:product:0BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: G
ReluReluBiasAdd:output:0*
T0*
_output_shapes

: X
IdentityIdentityRelu:activations:0^NoOp*
T0*
_output_shapes

: w
NoOpNoOp^BiasAdd/ReadVariableOp^MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*!
_input_shapes
: : : 20
BiasAdd/ReadVariableOpBiasAdd/ReadVariableOp2.
MatMul/ReadVariableOpMatMul/ReadVariableOp:F B

_output_shapes

: 
 
_user_specified_nameinputs
?	
?
@__inference_dense_1_layer_call_and_return_conditional_losses_855

inputs0
matmul_readvariableop_resource:-
biasadd_readvariableop_resource:
identity??BiasAdd/ReadVariableOp?MatMul/ReadVariableOpt
MatMul/ReadVariableOpReadVariableOpmatmul_readvariableop_resource*
_output_shapes

:*
dtype0`
MatMulMatMulinputsMatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: r
BiasAdd/ReadVariableOpReadVariableOpbiasadd_readvariableop_resource*
_output_shapes
:*
dtype0m
BiasAddBiasAddMatMul:product:0BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: M
SigmoidSigmoidBiasAdd:output:0*
T0*
_output_shapes

: Q
IdentityIdentitySigmoid:y:0^NoOp*
T0*
_output_shapes

: w
NoOpNoOp^BiasAdd/ReadVariableOp^MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*!
_input_shapes
: : : 20
BiasAdd/ReadVariableOpBiasAdd/ReadVariableOp2.
MatMul/ReadVariableOpMatMul/ReadVariableOp:F B

_output_shapes

: 
 
_user_specified_nameinputs
?
?
(__inference_sequential_layer_call_fn_946
input_1
unknown:
	unknown_0:
	unknown_1:
	unknown_2:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCallinput_1unknown	unknown_0	unknown_1	unknown_2*
Tin	
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *&
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *L
fGRE
C__inference_sequential_layer_call_and_return_conditional_losses_922f
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 22
StatefulPartitionedCallStatefulPartitionedCall:G C

_output_shapes

: 
!
_user_specified_name	input_1
?
?
!__inference_internal_grad_fn_1147
result_grads_0
result_grads_1
result_grads_2
result_grads_3
result_grads_4
result_grads_5
result_grads_6
result_grads_7

identity_4

identity_5

identity_6

identity_7M
IdentityIdentityresult_grads_0*
T0*
_output_shapes

:K

Identity_1Identityresult_grads_1*
T0*
_output_shapes
:O

Identity_2Identityresult_grads_2*
T0*
_output_shapes

:K

Identity_3Identityresult_grads_3*
T0*
_output_shapes
:?
	IdentityN	IdentityNresult_grads_0result_grads_1result_grads_2result_grads_3result_grads_0result_grads_1result_grads_2result_grads_3*
T

2**
_gradient_op_typeCustomGradient-1130*T
_output_shapesB
@::::::::S

Identity_4IdentityIdentityN:output:0*
T0*
_output_shapes

:O

Identity_5IdentityIdentityN:output:1*
T0*
_output_shapes
:S

Identity_6IdentityIdentityN:output:2*
T0*
_output_shapes

:O

Identity_7IdentityIdentityN:output:3*
T0*
_output_shapes
:"!

identity_4Identity_4:output:0"!

identity_5Identity_5:output:0"!

identity_6Identity_6:output:0"!

identity_7Identity_7:output:0*S
_input_shapesB
@:::::::::N J

_output_shapes

:
(
_user_specified_nameresult_grads_0:JF

_output_shapes
:
(
_user_specified_nameresult_grads_1:NJ

_output_shapes

:
(
_user_specified_nameresult_grads_2:JF

_output_shapes
:
(
_user_specified_nameresult_grads_3:NJ

_output_shapes

:
(
_user_specified_nameresult_grads_4:JF

_output_shapes
:
(
_user_specified_nameresult_grads_5:NJ

_output_shapes

:
(
_user_specified_nameresult_grads_6:JF

_output_shapes
:
(
_user_specified_nameresult_grads_7
?
?
__inference_predict_786
dataA
/sequential_dense_matmul_readvariableop_resource:>
0sequential_dense_biasadd_readvariableop_resource:C
1sequential_dense_1_matmul_readvariableop_resource:@
2sequential_dense_1_biasadd_readvariableop_resource:
identity??'sequential/dense/BiasAdd/ReadVariableOp?&sequential/dense/MatMul/ReadVariableOp?)sequential/dense_1/BiasAdd/ReadVariableOp?(sequential/dense_1/MatMul/ReadVariableOp?
&sequential/dense/MatMul/ReadVariableOpReadVariableOp/sequential_dense_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
sequential/dense/MatMulMatMuldata.sequential/dense/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
'sequential/dense/BiasAdd/ReadVariableOpReadVariableOp0sequential_dense_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
sequential/dense/BiasAddBiasAdd!sequential/dense/MatMul:product:0/sequential/dense/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: i
sequential/dense/ReluRelu!sequential/dense/BiasAdd:output:0*
T0*
_output_shapes

: ?
(sequential/dense_1/MatMul/ReadVariableOpReadVariableOp1sequential_dense_1_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
sequential/dense_1/MatMulMatMul#sequential/dense/Relu:activations:00sequential/dense_1/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
)sequential/dense_1/BiasAdd/ReadVariableOpReadVariableOp2sequential_dense_1_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
sequential/dense_1/BiasAddBiasAdd#sequential/dense_1/MatMul:product:01sequential/dense_1/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: s
sequential/dense_1/SigmoidSigmoid#sequential/dense_1/BiasAdd:output:0*
T0*
_output_shapes

: d
IdentityIdentitysequential/dense_1/Sigmoid:y:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp(^sequential/dense/BiasAdd/ReadVariableOp'^sequential/dense/MatMul/ReadVariableOp*^sequential/dense_1/BiasAdd/ReadVariableOp)^sequential/dense_1/MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2R
'sequential/dense/BiasAdd/ReadVariableOp'sequential/dense/BiasAdd/ReadVariableOp2P
&sequential/dense/MatMul/ReadVariableOp&sequential/dense/MatMul/ReadVariableOp2V
)sequential/dense_1/BiasAdd/ReadVariableOp)sequential/dense_1/BiasAdd/ReadVariableOp2T
(sequential/dense_1/MatMul/ReadVariableOp(sequential/dense_1/MatMul/ReadVariableOp:M I
'
_output_shapes
:?????????

_user_specified_namedata
?	
?
>__inference_dense_layer_call_and_return_conditional_losses_838

inputs0
matmul_readvariableop_resource:-
biasadd_readvariableop_resource:
identity??BiasAdd/ReadVariableOp?MatMul/ReadVariableOpt
MatMul/ReadVariableOpReadVariableOpmatmul_readvariableop_resource*
_output_shapes

:*
dtype0`
MatMulMatMulinputsMatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: r
BiasAdd/ReadVariableOpReadVariableOpbiasadd_readvariableop_resource*
_output_shapes
:*
dtype0m
BiasAddBiasAddMatMul:product:0BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: G
ReluReluBiasAdd:output:0*
T0*
_output_shapes

: X
IdentityIdentityRelu:activations:0^NoOp*
T0*
_output_shapes

: w
NoOpNoOp^BiasAdd/ReadVariableOp^MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*!
_input_shapes
: : : 20
BiasAdd/ReadVariableOpBiasAdd/ReadVariableOp2.
MatMul/ReadVariableOpMatMul/ReadVariableOp:F B

_output_shapes

: 
 
_user_specified_nameinputs
?
?
!__inference_signature_wrapper_801
data
unknown:
	unknown_0:
	unknown_1:
	unknown_2:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCalldataunknown	unknown_0	unknown_1	unknown_2*
Tin	
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *&
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? * 
fR
__inference_predict_786f
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*.
_input_shapes
:?????????: : : : 22
StatefulPartitionedCallStatefulPartitionedCall:M I
'
_output_shapes
:?????????

_user_specified_namedata
?J
?

 __inference__traced_restore_1248
file_prefix#
assignvariableop_variable: &
assignvariableop_1_adam_iter:	 (
assignvariableop_2_adam_beta_1: (
assignvariableop_3_adam_beta_2: '
assignvariableop_4_adam_decay: /
%assignvariableop_5_adam_learning_rate: 1
assignvariableop_6_dense_kernel:+
assignvariableop_7_dense_bias:3
!assignvariableop_8_dense_1_kernel:-
assignvariableop_9_dense_1_bias:9
'assignvariableop_10_adam_dense_kernel_m:3
%assignvariableop_11_adam_dense_bias_m:;
)assignvariableop_12_adam_dense_1_kernel_m:5
'assignvariableop_13_adam_dense_1_bias_m:9
'assignvariableop_14_adam_dense_kernel_v:3
%assignvariableop_15_adam_dense_bias_v:;
)assignvariableop_16_adam_dense_1_kernel_v:5
'assignvariableop_17_adam_dense_1_bias_v:
identity_19??AssignVariableOp?AssignVariableOp_1?AssignVariableOp_10?AssignVariableOp_11?AssignVariableOp_12?AssignVariableOp_13?AssignVariableOp_14?AssignVariableOp_15?AssignVariableOp_16?AssignVariableOp_17?AssignVariableOp_2?AssignVariableOp_3?AssignVariableOp_4?AssignVariableOp_5?AssignVariableOp_6?AssignVariableOp_7?AssignVariableOp_8?AssignVariableOp_9?

RestoreV2/tensor_namesConst"/device:CPU:0*
_output_shapes
:*
dtype0*?

value?
B?
B'_global_step/.ATTRIBUTES/VARIABLE_VALUEB*_optimizer/iter/.ATTRIBUTES/VARIABLE_VALUEB,_optimizer/beta_1/.ATTRIBUTES/VARIABLE_VALUEB,_optimizer/beta_2/.ATTRIBUTES/VARIABLE_VALUEB+_optimizer/decay/.ATTRIBUTES/VARIABLE_VALUEB3_optimizer/learning_rate/.ATTRIBUTES/VARIABLE_VALUEB=_model/layer_with_weights-0/kernel/.ATTRIBUTES/VARIABLE_VALUEB;_model/layer_with_weights-0/bias/.ATTRIBUTES/VARIABLE_VALUEB=_model/layer_with_weights-1/kernel/.ATTRIBUTES/VARIABLE_VALUEB;_model/layer_with_weights-1/bias/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-0/kernel/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-0/bias/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-1/kernel/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-1/bias/.OPTIMIZER_SLOT/_optimizer/m/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-0/kernel/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-0/bias/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEBZ_model/layer_with_weights-1/kernel/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEBX_model/layer_with_weights-1/bias/.OPTIMIZER_SLOT/_optimizer/v/.ATTRIBUTES/VARIABLE_VALUEB_CHECKPOINTABLE_OBJECT_GRAPH?
RestoreV2/shape_and_slicesConst"/device:CPU:0*
_output_shapes
:*
dtype0*9
value0B.B B B B B B B B B B B B B B B B B B B ?
	RestoreV2	RestoreV2file_prefixRestoreV2/tensor_names:output:0#RestoreV2/shape_and_slices:output:0"/device:CPU:0*`
_output_shapesN
L:::::::::::::::::::*!
dtypes
2	[
IdentityIdentityRestoreV2:tensors:0"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOpAssignVariableOpassignvariableop_variableIdentity:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_1IdentityRestoreV2:tensors:1"/device:CPU:0*
T0	*
_output_shapes
:?
AssignVariableOp_1AssignVariableOpassignvariableop_1_adam_iterIdentity_1:output:0"/device:CPU:0*
_output_shapes
 *
dtype0	]

Identity_2IdentityRestoreV2:tensors:2"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_2AssignVariableOpassignvariableop_2_adam_beta_1Identity_2:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_3IdentityRestoreV2:tensors:3"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_3AssignVariableOpassignvariableop_3_adam_beta_2Identity_3:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_4IdentityRestoreV2:tensors:4"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_4AssignVariableOpassignvariableop_4_adam_decayIdentity_4:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_5IdentityRestoreV2:tensors:5"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_5AssignVariableOp%assignvariableop_5_adam_learning_rateIdentity_5:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_6IdentityRestoreV2:tensors:6"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_6AssignVariableOpassignvariableop_6_dense_kernelIdentity_6:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_7IdentityRestoreV2:tensors:7"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_7AssignVariableOpassignvariableop_7_dense_biasIdentity_7:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_8IdentityRestoreV2:tensors:8"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_8AssignVariableOp!assignvariableop_8_dense_1_kernelIdentity_8:output:0"/device:CPU:0*
_output_shapes
 *
dtype0]

Identity_9IdentityRestoreV2:tensors:9"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_9AssignVariableOpassignvariableop_9_dense_1_biasIdentity_9:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_10IdentityRestoreV2:tensors:10"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_10AssignVariableOp'assignvariableop_10_adam_dense_kernel_mIdentity_10:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_11IdentityRestoreV2:tensors:11"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_11AssignVariableOp%assignvariableop_11_adam_dense_bias_mIdentity_11:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_12IdentityRestoreV2:tensors:12"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_12AssignVariableOp)assignvariableop_12_adam_dense_1_kernel_mIdentity_12:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_13IdentityRestoreV2:tensors:13"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_13AssignVariableOp'assignvariableop_13_adam_dense_1_bias_mIdentity_13:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_14IdentityRestoreV2:tensors:14"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_14AssignVariableOp'assignvariableop_14_adam_dense_kernel_vIdentity_14:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_15IdentityRestoreV2:tensors:15"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_15AssignVariableOp%assignvariableop_15_adam_dense_bias_vIdentity_15:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_16IdentityRestoreV2:tensors:16"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_16AssignVariableOp)assignvariableop_16_adam_dense_1_kernel_vIdentity_16:output:0"/device:CPU:0*
_output_shapes
 *
dtype0_
Identity_17IdentityRestoreV2:tensors:17"/device:CPU:0*
T0*
_output_shapes
:?
AssignVariableOp_17AssignVariableOp'assignvariableop_17_adam_dense_1_bias_vIdentity_17:output:0"/device:CPU:0*
_output_shapes
 *
dtype01
NoOpNoOp"/device:CPU:0*
_output_shapes
 ?
Identity_18Identityfile_prefix^AssignVariableOp^AssignVariableOp_1^AssignVariableOp_10^AssignVariableOp_11^AssignVariableOp_12^AssignVariableOp_13^AssignVariableOp_14^AssignVariableOp_15^AssignVariableOp_16^AssignVariableOp_17^AssignVariableOp_2^AssignVariableOp_3^AssignVariableOp_4^AssignVariableOp_5^AssignVariableOp_6^AssignVariableOp_7^AssignVariableOp_8^AssignVariableOp_9^NoOp"/device:CPU:0*
T0*
_output_shapes
: W
Identity_19IdentityIdentity_18:output:0^NoOp_1*
T0*
_output_shapes
: ?
NoOp_1NoOp^AssignVariableOp^AssignVariableOp_1^AssignVariableOp_10^AssignVariableOp_11^AssignVariableOp_12^AssignVariableOp_13^AssignVariableOp_14^AssignVariableOp_15^AssignVariableOp_16^AssignVariableOp_17^AssignVariableOp_2^AssignVariableOp_3^AssignVariableOp_4^AssignVariableOp_5^AssignVariableOp_6^AssignVariableOp_7^AssignVariableOp_8^AssignVariableOp_9*"
_acd_function_control_output(*
_output_shapes
 "#
identity_19Identity_19:output:0*9
_input_shapes(
&: : : : : : : : : : : : : : : : : : : 2$
AssignVariableOpAssignVariableOp2(
AssignVariableOp_1AssignVariableOp_12*
AssignVariableOp_10AssignVariableOp_102*
AssignVariableOp_11AssignVariableOp_112*
AssignVariableOp_12AssignVariableOp_122*
AssignVariableOp_13AssignVariableOp_132*
AssignVariableOp_14AssignVariableOp_142*
AssignVariableOp_15AssignVariableOp_152*
AssignVariableOp_16AssignVariableOp_162*
AssignVariableOp_17AssignVariableOp_172(
AssignVariableOp_2AssignVariableOp_22(
AssignVariableOp_3AssignVariableOp_32(
AssignVariableOp_4AssignVariableOp_42(
AssignVariableOp_5AssignVariableOp_52(
AssignVariableOp_6AssignVariableOp_62(
AssignVariableOp_7AssignVariableOp_72(
AssignVariableOp_8AssignVariableOp_82(
AssignVariableOp_9AssignVariableOp_9:C ?

_output_shapes
: 
%
_user_specified_namefile_prefix
?
?
(__inference_sequential_layer_call_fn_873
input_1
unknown:
	unknown_0:
	unknown_1:
	unknown_2:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCallinput_1unknown	unknown_0	unknown_1	unknown_2*
Tin	
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *&
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *L
fGRE
C__inference_sequential_layer_call_and_return_conditional_losses_862f
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 22
StatefulPartitionedCallStatefulPartitionedCall:G C

_output_shapes

: 
!
_user_specified_name	input_1
?
?
!__inference_signature_wrapper_768
data

labels
unknown: 
	unknown_0:
	unknown_1:
	unknown_2:
	unknown_3:
	unknown_4: 
	unknown_5:	 
	unknown_6: 
	unknown_7: 
	unknown_8:
	unknown_9:

unknown_10:

unknown_11:

unknown_12:

unknown_13:

unknown_14:

unknown_15:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCalldatalabelsunknown	unknown_0	unknown_1	unknown_2	unknown_3	unknown_4	unknown_5	unknown_6	unknown_7	unknown_8	unknown_9
unknown_10
unknown_11
unknown_12
unknown_13
unknown_14
unknown_15*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes
: *%
_read_only_resource_inputs
	
*-
config_proto

CPU

GPU 2J 8? *
fR
__inference_learn_726^
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes
: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*W
_input_shapesF
D:?????????:?????????: : : : : : : : : : : : : : : : : 22
StatefulPartitionedCallStatefulPartitionedCall:M I
'
_output_shapes
:?????????

_user_specified_namedata:KG
#
_output_shapes
:?????????
 
_user_specified_namelabels
?
?
C__inference_sequential_layer_call_and_return_conditional_losses_922

inputs
	dense_911:
	dense_913:
dense_1_916:
dense_1_918:
identity??dense/StatefulPartitionedCall?dense_1/StatefulPartitionedCall?
dense/StatefulPartitionedCallStatefulPartitionedCallinputs	dense_911	dense_913*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *G
fBR@
>__inference_dense_layer_call_and_return_conditional_losses_838?
dense_1/StatefulPartitionedCallStatefulPartitionedCall&dense/StatefulPartitionedCall:output:0dense_1_916dense_1_918*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *I
fDRB
@__inference_dense_1_layer_call_and_return_conditional_losses_855n
IdentityIdentity(dense_1/StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp^dense/StatefulPartitionedCall ^dense_1/StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2>
dense/StatefulPartitionedCalldense/StatefulPartitionedCall2B
dense_1/StatefulPartitionedCalldense_1/StatefulPartitionedCall:F B

_output_shapes

: 
 
_user_specified_nameinputs
?
?
$__inference_dense_layer_call_fn_1045

inputs
unknown:
	unknown_0:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCallinputsunknown	unknown_0*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *G
fBR@
>__inference_dense_layer_call_and_return_conditional_losses_838f
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*!
_input_shapes
: : : 22
StatefulPartitionedCallStatefulPartitionedCall:F B

_output_shapes

: 
 
_user_specified_nameinputs
?
?
D__inference_sequential_layer_call_and_return_conditional_losses_1018

inputs6
$dense_matmul_readvariableop_resource:3
%dense_biasadd_readvariableop_resource:8
&dense_1_matmul_readvariableop_resource:5
'dense_1_biasadd_readvariableop_resource:
identity??dense/BiasAdd/ReadVariableOp?dense/MatMul/ReadVariableOp?dense_1/BiasAdd/ReadVariableOp?dense_1/MatMul/ReadVariableOp?
dense/MatMul/ReadVariableOpReadVariableOp$dense_matmul_readvariableop_resource*
_output_shapes

:*
dtype0l
dense/MatMulMatMulinputs#dense/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ~
dense/BiasAdd/ReadVariableOpReadVariableOp%dense_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0
dense/BiasAddBiasAdddense/MatMul:product:0$dense/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: S

dense/ReluReludense/BiasAdd:output:0*
T0*
_output_shapes

: ?
dense_1/MatMul/ReadVariableOpReadVariableOp&dense_1_matmul_readvariableop_resource*
_output_shapes

:*
dtype0?
dense_1/MatMulMatMuldense/Relu:activations:0%dense_1/MatMul/ReadVariableOp:value:0*
T0*
_output_shapes

: ?
dense_1/BiasAdd/ReadVariableOpReadVariableOp'dense_1_biasadd_readvariableop_resource*
_output_shapes
:*
dtype0?
dense_1/BiasAddBiasAdddense_1/MatMul:product:0&dense_1/BiasAdd/ReadVariableOp:value:0*
T0*
_output_shapes

: ]
dense_1/SigmoidSigmoiddense_1/BiasAdd:output:0*
T0*
_output_shapes

: Y
IdentityIdentitydense_1/Sigmoid:y:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp^dense/BiasAdd/ReadVariableOp^dense/MatMul/ReadVariableOp^dense_1/BiasAdd/ReadVariableOp^dense_1/MatMul/ReadVariableOp*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2<
dense/BiasAdd/ReadVariableOpdense/BiasAdd/ReadVariableOp2:
dense/MatMul/ReadVariableOpdense/MatMul/ReadVariableOp2@
dense_1/BiasAdd/ReadVariableOpdense_1/BiasAdd/ReadVariableOp2>
dense_1/MatMul/ReadVariableOpdense_1/MatMul/ReadVariableOp:F B

_output_shapes

: 
 
_user_specified_nameinputs
?
?
C__inference_sequential_layer_call_and_return_conditional_losses_960
input_1
	dense_949:
	dense_951:
dense_1_954:
dense_1_956:
identity??dense/StatefulPartitionedCall?dense_1/StatefulPartitionedCall?
dense/StatefulPartitionedCallStatefulPartitionedCallinput_1	dense_949	dense_951*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *G
fBR@
>__inference_dense_layer_call_and_return_conditional_losses_838?
dense_1/StatefulPartitionedCallStatefulPartitionedCall&dense/StatefulPartitionedCall:output:0dense_1_954dense_1_956*
Tin
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *$
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *I
fDRB
@__inference_dense_1_layer_call_and_return_conditional_losses_855n
IdentityIdentity(dense_1/StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: ?
NoOpNoOp^dense/StatefulPartitionedCall ^dense_1/StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 2>
dense/StatefulPartitionedCalldense/StatefulPartitionedCall2B
dense_1/StatefulPartitionedCalldense_1/StatefulPartitionedCall:G C

_output_shapes

: 
!
_user_specified_name	input_1
?
?
(__inference_sequential_layer_call_fn_987

inputs
unknown:
	unknown_0:
	unknown_1:
	unknown_2:
identity??StatefulPartitionedCall?
StatefulPartitionedCallStatefulPartitionedCallinputsunknown	unknown_0	unknown_1	unknown_2*
Tin	
2*
Tout
2*
_collective_manager_ids
 *
_output_shapes

: *&
_read_only_resource_inputs
*-
config_proto

CPU

GPU 2J 8? *L
fGRE
C__inference_sequential_layer_call_and_return_conditional_losses_862f
IdentityIdentity StatefulPartitionedCall:output:0^NoOp*
T0*
_output_shapes

: `
NoOpNoOp^StatefulPartitionedCall*"
_acd_function_control_output(*
_output_shapes
 "
identityIdentity:output:0*(
_construction_contextkEagerRuntime*%
_input_shapes
: : : : : 22
StatefulPartitionedCallStatefulPartitionedCall:F B

_output_shapes

: 
 
_user_specified_nameinputs7
!__inference_internal_grad_fn_1147CustomGradient-699"?L
saver_filename:0StatefulPartitionedCall_2:0StatefulPartitionedCall_38"
saved_model_main_op

NoOp*>
__saved_model_init_op%#
__saved_model_init_op

NoOp*?
learn?
+
data#
learn_data:0?????????
+
labels!
learn_labels:0?????????'
loss
StatefulPartitionedCall:0 tensorflow/serving/predict*?
predict?
-
data%
predict_data:0?????????8
predictions)
StatefulPartitionedCall_1:0 tensorflow/serving/predict:?F
t

_model
_global_step

_optimizer

signatures
	3learn
4predict"
_generic_user_object
?
layer_with_weights-0
layer-0
layer_with_weights-1
layer-1
	variables
trainable_variables
	regularization_losses

	keras_api
5__call__
*6&call_and_return_all_conditional_losses
7_default_save_signature"
_tf_keras_sequential
: 2Variable
?
iter

beta_1

beta_2
	decay
learning_ratem+m,m-m.v/v0v1v2"
	optimizer
/
	8learn
9predict"
signature_map
?

kernel
bias
	variables
trainable_variables
regularization_losses
	keras_api
:__call__
*;&call_and_return_all_conditional_losses"
_tf_keras_layer
?

kernel
bias
	variables
trainable_variables
regularization_losses
	keras_api
<__call__
*=&call_and_return_all_conditional_losses"
_tf_keras_layer
<
0
1
2
3"
trackable_list_wrapper
<
0
1
2
3"
trackable_list_wrapper
 "
trackable_list_wrapper
?
non_trainable_variables

layers
metrics
layer_regularization_losses
 layer_metrics
	variables
trainable_variables
	regularization_losses
5__call__
7_default_save_signature
*6&call_and_return_all_conditional_losses
&6"call_and_return_conditional_losses"
_generic_user_object
:	 (2	Adam/iter
: (2Adam/beta_1
: (2Adam/beta_2
: (2
Adam/decay
: (2Adam/learning_rate
:2dense/kernel
:2
dense/bias
.
0
1"
trackable_list_wrapper
.
0
1"
trackable_list_wrapper
 "
trackable_list_wrapper
?
!non_trainable_variables

"layers
#metrics
$layer_regularization_losses
%layer_metrics
	variables
trainable_variables
regularization_losses
:__call__
*;&call_and_return_all_conditional_losses
&;"call_and_return_conditional_losses"
_generic_user_object
 :2dense_1/kernel
:2dense_1/bias
.
0
1"
trackable_list_wrapper
.
0
1"
trackable_list_wrapper
 "
trackable_list_wrapper
?
&non_trainable_variables

'layers
(metrics
)layer_regularization_losses
*layer_metrics
	variables
trainable_variables
regularization_losses
<__call__
*=&call_and_return_all_conditional_losses
&="call_and_return_conditional_losses"
_generic_user_object
 "
trackable_list_wrapper
.
0
1"
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_dict_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_dict_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_list_wrapper
 "
trackable_dict_wrapper
#:!2Adam/dense/kernel/m
:2Adam/dense/bias/m
%:#2Adam/dense_1/kernel/m
:2Adam/dense_1/bias/m
#:!2Adam/dense/kernel/v
:2Adam/dense/bias/v
%:#2Adam/dense_1/kernel/v
:2Adam/dense_1/bias/v
?2?
__inference_learn_726?
???
FullArgSpec%
args?
jself
jdata
jlabels
varargs
 
varkw
 
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *3?0
??????????
??????????
?2?
__inference_predict_786?
???
FullArgSpec
args?
jself
jdata
varargs
 
varkw
 
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *?
??????????
?2?
(__inference_sequential_layer_call_fn_873
(__inference_sequential_layer_call_fn_987
)__inference_sequential_layer_call_fn_1000
(__inference_sequential_layer_call_fn_946?
???
FullArgSpec1
args)?&
jself
jinputs

jtraining
jmask
varargs
 
varkw
 
defaults?
p 

 

kwonlyargs? 
kwonlydefaults? 
annotations? *
 
?2?
D__inference_sequential_layer_call_and_return_conditional_losses_1018
D__inference_sequential_layer_call_and_return_conditional_losses_1036
C__inference_sequential_layer_call_and_return_conditional_losses_960
C__inference_sequential_layer_call_and_return_conditional_losses_974?
???
FullArgSpec1
args)?&
jself
jinputs

jtraining
jmask
varargs
 
varkw
 
defaults?
p 

 

kwonlyargs? 
kwonlydefaults? 
annotations? *
 
?B?
__inference__wrapped_model_820input_1"?
???
FullArgSpec
args? 
varargsjargs
varkwjkwargs
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *
 
?B?
!__inference_signature_wrapper_768datalabels"?
???
FullArgSpec
args? 
varargs
 
varkwjkwargs
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *
 
?B?
!__inference_signature_wrapper_801data"?
???
FullArgSpec
args? 
varargs
 
varkwjkwargs
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *
 
?2?
$__inference_dense_layer_call_fn_1045?
???
FullArgSpec
args?
jself
jinputs
varargs
 
varkw
 
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *
 
?2?
?__inference_dense_layer_call_and_return_conditional_losses_1056?
???
FullArgSpec
args?
jself
jinputs
varargs
 
varkw
 
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *
 
?2?
&__inference_dense_1_layer_call_fn_1065?
???
FullArgSpec
args?
jself
jinputs
varargs
 
varkw
 
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *
 
?2?
A__inference_dense_1_layer_call_and_return_conditional_losses_1076?
???
FullArgSpec
args?
jself
jinputs
varargs
 
varkw
 
defaults
 

kwonlyargs? 
kwonlydefaults
 
annotations? *
 {
__inference__wrapped_model_820Y'?$
?
?
input_1 
? "(?%
#
dense_1?
dense_1 ?
A__inference_dense_1_layer_call_and_return_conditional_losses_1076J&?#
?
?
inputs 
? "?
?
0 
? g
&__inference_dense_1_layer_call_fn_1065=&?#
?
?
inputs 
? "? ?
?__inference_dense_layer_call_and_return_conditional_losses_1056J&?#
?
?
inputs 
? "?
?
0 
? e
$__inference_dense_layer_call_fn_1045=&?#
?
?
inputs 
? "? ?
!__inference_internal_grad_fn_1147????
???

 
?
result_grads_0
?
result_grads_1
?
result_grads_2
?
result_grads_3
?
result_grads_4
?
result_grads_5
?
result_grads_6
?
result_grads_7
? "[?X

 

 

 

 
?
4
?
5
?
6
?
7?
__inference_learn_726|+/,0-1.2K?H
A?>
?
data?????????
?
labels?????????
? "?

loss?

loss ?
__inference_predict_786g-?*
#? 
?
data?????????
? "0?-
+
predictions?
predictions ?
D__inference_sequential_layer_call_and_return_conditional_losses_1018T.?+
$?!
?
inputs 
p 

 
? "?
?
0 
? ?
D__inference_sequential_layer_call_and_return_conditional_losses_1036T.?+
$?!
?
inputs 
p

 
? "?
?
0 
? ?
C__inference_sequential_layer_call_and_return_conditional_losses_960U/?,
%?"
?
input_1 
p 

 
? "?
?
0 
? ?
C__inference_sequential_layer_call_and_return_conditional_losses_974U/?,
%?"
?
input_1 
p

 
? "?
?
0 
? t
)__inference_sequential_layer_call_fn_1000G.?+
$?!
?
inputs 
p

 
? "? t
(__inference_sequential_layer_call_fn_873H/?,
%?"
?
input_1 
p 

 
? "? t
(__inference_sequential_layer_call_fn_946H/?,
%?"
?
input_1 
p

 
? "? s
(__inference_sequential_layer_call_fn_987G.?+
$?!
?
inputs 
p 

 
? "? ?
!__inference_signature_wrapper_768?+/,0-1.2]?Z
? 
S?P
&
data?
data?????????
&
labels?
labels?????????"?

loss?

loss ?
!__inference_signature_wrapper_801o5?2
? 
+?(
&
data?
data?????????"0?-
+
predictions?
predictions 