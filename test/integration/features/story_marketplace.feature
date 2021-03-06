@story_marketplace
Feature: 
    Install OpenShift operator from OperatorHub and use it to manage
    admin tasks.

    Scenario Outline: Start CRC and login to cluster
        Given executing "crc setup" succeeds
        When starting CRC with default bundle and hypervisor "<vm-driver>" succeeds
        Then stdout should contain "Started the OpenShift cluster"
        And executing "eval $(crc oc-env)" succeeds
        When with up to "4" retries with wait period of "2m" command "crc status" output matches ".*Running \(v\d+\.\d+\.\d+.*\).*"
        Then login to the oc cluster succeeds

    @darwin
        Examples:
            | vm-driver  |
            | hyperkit   |

    @linux
        Examples:
            | vm-driver |
            | libvirt   |

    @windows
    Scenario: Start CRC on Windows
        Given executing "crc setup" succeeds
        When starting CRC with default bundle and nameserver "10.75.5.25" succeeds
        Then stdout should contain "Started the OpenShift cluster"
        And executing "crc oc-env | Invoke-Expression" succeeds
        When with up to "4" retries with wait period of "2m" command "crc status" output matches ".*Running \(v\d+\.\d+\.\d+.*\).*"
        Then login to the oc cluster succeeds

    @darwin @linux @windows
    Scenario: Install new operator
        Given executing "oc apply -f etcdop-csc.yaml" succeeds
        When executing "oc apply -f etcdop-sub.yaml" succeeds
        # check if cluster operator is running
        Then with up to "20" retries with wait period of "30s" command "oc get csv" output matches ".*etcdoperator\.(.*)Succeeded$"
        
    @darwin @linux @windows
    Scenario: Scale up
        # start cluster with 3 pods
        When executing "oc apply -f etcd-cluster3.yaml" succeeds
        Then with up to "10" retries with wait period of "30s" command "oc get pods" output matches "(?s)(.*example-[a-z0-9]* *1/1 *Running.*){3}"
        # reconfigure cluster to 5 pods
        When executing "oc apply -f etcd-cluster5.yaml" succeeds
        Then with up to "10" retries with wait period of "30s" command "oc get pods" output matches "(?s)(.*example-[a-z0-9]* *1/1 *Running.*){5}"
    
    @darwin @linux @windows
    Scenario: Failover
        # simulate failure of 1 pod, check that it was replaced
        When executing "POD=$(oc get pod -o jsonpath="{.items[0].metadata.name}")" succeeds
        And executing "echo $POD" succeeds
        And executing "oc delete pod $POD --now" succeeds
        Then stdout should match "^pod(.*)deleted$"
        # after a while 5 pods should be up & running again
        And with up to "10" retries with wait period of "30s" command "oc get pods" output matches "(?s)(.*example-[a-z0-9]* *1/1 *Running.*){5}"
        # but the deleted pod should not be up, it was replaced
        But executing "oc get pods $POD" fails
        And stderr matches "(.*)pods (.*) not found$"

    @darwin @linux @windows
    Scenario: Scale down
        # scale back down to 3 pods
        When executing "oc apply -f etcd-cluster3.yaml" succeeds
        Then with up to "10" retries with wait period of "30s" command "oc get pods" output matches "(?s)(.*example-[a-z0-9]* *1/1 *Running.*){3}"
        But with up to "10" retries with wait period of "30s" command "oc get pods" output does not match "(?s)(.*example-[a-z0-9]* *1/1 *Running.*){4}"
        
    @darwin @linux @windows
    Scenario: Clean up
        When executing "crc stop -f" succeeds
        Then stdout should match "(.*)[Ss]topped the OpenShift cluster"
        When executing "crc delete -f" succeeds
        Then stdout should contain "Deleted the OpenShift cluster"
