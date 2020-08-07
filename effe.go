// Package effe contains directives for Effe code generation.
// For an overview of working with Effe, see the user guide at
// https://github.com/GettEngineering/effe/blob/master/docs/docs/GettingStarted.md
//
// The directives in this package are used as input to the Effe code generation
// tool. The entry point of Effe's analysis are injector functions: function
// templates denoted by only containing a call to BuildFlow. The arguments to BuildFlow
// describes a set of steps and the Effe code generation tool builds
// calls of steps according to strategy.
package effe

// Based type for declaring steps
type StepFunc interface{}

// Type for declaring business logic for specific case
type CaseKey interface{}

func panicDSLMethodNotFound() {
	panic("implementation not generated, run effe before")
}

// Build is placed in the body of an inhector function template to declare the steps.
func BuildFlow(funcs ...StepFunc) interface{} { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}

// A Step declares a function which will be executed in this place.
// The function should have the following format:
//
// Format:
//
//          func step1(dep1 FirstDependency, dep2 SecondDependency) func(req *httpRequest) error {
//              return func(req *httpRequest) error {
//                  return dep1.Send(dep2, req)
//              }
//          }
//
// Function arguments should be dependencies for your step. It is necessary
// to separate dependencies between steps. You don't need to create big service objects.
// Effe code generation tool calculates all dependencies for your flow and Effe generates
// a service object with dependecies automatically.
// Function return value should be the function which executes here.
// Also, you can call another business flow here. It helps to split and reuse existing business logic.
//
// Examples:
//
//          func BuildMyFirstBusinessFlow(){
//              effe.BuldFlow(
//                  effe.Step(step1),
//              )
//          }
//
//          func BuildMySecondBusinessFlow(){
//              effe.BuildFlow(
//                  effe.Step(BuildMyFirstBusinessFlow),
//              )
//          }
func Step(fn interface{}) StepFunc { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}

// A Failure works the same way as Step, but with one exception:
// a function executes only if one of steps returns an error.
func Failure(fn interface{}) StepFunc { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}

// Wrap helps to declare a block of steps which is executed in the following sequence.
// If we declare a function with Before directive, then it will be executed first.
// After that, all functions declared with Step directive will be executed.
// Lastly a function declared with Success directive will be executed.
// Also, we can declare an error handler here with Failure directive.
//
// Examples:
//
//          func BuildMyBusinessFlow() {
//              effe.BuildFlow(
//                  effe.Step(step1),
//                  effe.Step(step2),
//                  effe.Wrap(effe.Before(lock), effe.Success(unlock), effe.Failure(catchErrAndUnlock),
//                      effe.Step(step3),
//                      effe.Step(step4),
//                  ),
//              )
//          }
func Wrap(beforeFunc StepFunc, afterFunc StepFunc, steps ...StepFunc) StepFunc { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}

// This directive helps to declare a function which executes after
// other steps in a directive Wrap and if not one step returned an error.
// This directive can be used only in Wrap.
func Success(fn interface{}) StepFunc { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}

// This directive helps to declare a function which executes before
// other steps in a directive Wrap.
// This directive can be used only in Wrap.
func Before(fn interface{}) StepFunc { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}

// Decision directive helps to organize branching of our business logic.
// First argument can be a type or a field from the type.
// Golang doesn't provide an opportunity to pass types to function argument.
// For that you need to create an empty object with your type.
// If you want branching your logic by field from type you can declare it.
// For that you need to create an empty object with your type and get field from it.
// Other arguments declare handlers for every case. Every case can be declared with Case directive.
// Failure directive works here too. With Failure you can declare an error handler for case.
//
// Examples:
//
//      func BuildMyBusinessFlow(){
//          effe.BuildFlow(
//              effe.Step(tryLock),
//              effe.Decision(new(Lock),
//                  effe.Case(true, effe.Step(step1)),
//                  effe.Case(false, effe.Step(step2)),
//              ),
//          )
//      }
//
//      func BuildMyBusinessFlow2(){
//          effe.BuildFlow(
//              effe.Step(tryLock),
//                  effe.Decision(Locker{}.Lock,
//                  effe.Case(true, effe.Step(step1)),
//                  effe.Case(false, effe.Step(step2)),
//              ),
//          )
//      }
func Decision(tag interface{}, cases ...StepFunc) StepFunc { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}

// This directive can be used only in Decision.
// Case declares the steps which will execute.
func Case(key CaseKey, funcs ...StepFunc) StepFunc { //nolint:unparam
	panicDSLMethodNotFound()
	return nil
}
